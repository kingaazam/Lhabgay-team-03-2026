package controllers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"lhabgay/backend/database"
	"lhabgay/backend/models"
	"lhabgay/backend/utils"

	"github.com/gorilla/mux"
)

// UploadBook saves uploaded files and creates a books table record.
func UploadBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		utils.Error(w, http.StatusBadRequest, "could not read upload form")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	author := strings.TrimSpace(r.FormValue("author"))
	description := strings.TrimSpace(r.FormValue("description"))
	category := strings.TrimSpace(r.FormValue("category"))
	if title == "" || author == "" || description == "" || category == "" {
		utils.Error(w, http.StatusBadRequest, "title, author, description and category are required")
		return
	}

	coverPath, err := saveUploadedFile(r, "cover_image", "image")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	bookPath, err := saveUploadedFile(r, "book_file", "book")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec(
		`INSERT INTO books (title, author, description, category, cover_image_path, book_file_path)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		title,
		author,
		description,
		category,
		coverPath,
		bookPath,
	)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "could not save book")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "book uploaded successfully"})
}

// GetBooks returns all books from newest to oldest.
func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, title, author, description, category, cover_image_path, book_file_path, created_at
		 FROM books
		 ORDER BY created_at DESC, id DESC`,
	)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "could not load books")
		return
	}
	defer rows.Close()

	books := make([]models.Book, 0)
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Description,
			&book.Category,
			&book.CoverImagePath,
			&book.BookFilePath,
			&book.CreatedAt,
		); err != nil {
			utils.Error(w, http.StatusInternalServerError, "could not read books")
			return
		}
		books = append(books, book)
	}

	utils.JSON(w, http.StatusOK, books)
}

// GetBook returns one book by URL id.
func GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid book id")
		return
	}

	var book models.Book
	err = database.DB.QueryRow(
		`SELECT id, title, author, description, category, cover_image_path, book_file_path, created_at
		 FROM books
		 WHERE id = $1`,
		id,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Description,
		&book.Category,
		&book.CoverImagePath,
		&book.BookFilePath,
		&book.CreatedAt,
	)
	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "could not load book")
		return
	}

	utils.JSON(w, http.StatusOK, book)
}

// DeleteBook removes one book record and its uploaded files. Admin middleware protects this route.
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid book id")
		return
	}

	var coverPath, bookPath string
	err = database.DB.QueryRow(
		"SELECT cover_image_path, book_file_path FROM books WHERE id = $1",
		id,
	).Scan(&coverPath, &bookPath)
	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "could not find book")
		return
	}

	result, err := database.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "could not delete book")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		utils.Error(w, http.StatusNotFound, "book not found")
		return
	}

	removeUploadedFile(coverPath)
	removeUploadedFile(bookPath)

	utils.JSON(w, http.StatusOK, map[string]string{"message": "book deleted successfully"})
}

func saveUploadedFile(r *http.Request, formName, folder string) (string, error) {
	file, header, err := r.FormFile(formName)
	if err != nil {
		return "", fmt.Errorf("%s is required", formName)
	}
	defer file.Close()

	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", fmt.Errorf("could not create %s folder", folder)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	base := strings.TrimSuffix(filepath.Base(header.Filename), ext)
	base = safeFileName(base)
	if base == "" {
		base = "upload"
	}

	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), base, ext)
	relativePath := filepath.ToSlash(filepath.Join(folder, fileName))
	destination, err := os.Create(relativePath)
	if err != nil {
		return "", fmt.Errorf("could not save %s", formName)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return "", fmt.Errorf("could not write %s", formName)
	}

	return relativePath, nil
}

func safeFileName(name string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '-' || r == '_':
			builder.WriteRune(r)
		case r == ' ':
			builder.WriteRune('_')
		}
	}
	return builder.String()
}

func removeUploadedFile(path string) {
	path = filepath.Clean(path)
	if strings.HasPrefix(path, "image"+string(os.PathSeparator)) || strings.HasPrefix(path, "book"+string(os.PathSeparator)) {
		_ = os.Remove(path)
	}
}

// ServeFile handles HTTP requests to serve uploaded static files (images, PDFs, etc.)
func ServeFile(w http.ResponseWriter, r *http.Request) {
	// Get the filepath variable from the mux URL router
	vars := mux.Vars(r)
	filePath := vars["filepath"]

	// Define the base directory where your files are stored locally
	// Change "./uploads" to "./image" or whichever folder your backend uses to save uploads
	baseDir := "./image"

	// Securely join paths to prevent directory traversal attacks (e.g., ../../etc/passwd)
	finalPath := filepath.Join(baseDir, filePath)

	// Serve the static file
	http.ServeFile(w, r, finalPath)
}
