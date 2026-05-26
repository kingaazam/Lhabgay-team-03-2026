package controllers

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gorilla/mux"
)

func ServeFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filepath"]
	filePath := filepath.Join("../book", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}

func UploadBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	category := r.FormValue("category")
	description := r.FormValue("description")

	// 1. Get the PDF Book File cleanly
	file, handler, err := r.FormFile("book_file")
	if err != nil {
		file, handler, err = r.FormFile("bookFile")
		if err != nil {
			file, handler, err = r.FormFile("file")
			if err != nil {
				http.Error(w, `{"error": "PDF file is required"}`, http.StatusBadRequest)
				return
			}
		}
	}
	defer file.Close()

	// 2. Get the Cover Image File cleanly
	coverFile, coverHandler, err := r.FormFile("cover_image")
	if err != nil {
		coverFile, coverHandler, err = r.FormFile("coverImage")
		if err != nil {
			// If cover image is optional, you can remove this error block
			http.Error(w, `{"error": "Cover image is required"}`, http.StatusBadRequest)
			return
		}
	}
	defer coverFile.Close()

	uploadDir := "../book"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Save Book PDF to disk
	filePath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save file content")
		return
	}

	// Save Cover Image to disk
	coverPath := filepath.Join(uploadDir, coverHandler.Filename)
	dstCover, err := os.Create(coverPath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save cover image")
		return
	}
	defer dstCover.Close()
	if _, err = io.Copy(dstCover, coverFile); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save cover content")
		return
	}

	// 3. Updated Database Execution (Includes both path fields)
	query := "INSERT INTO books (title, author, category, description, book_file_path, cover_image) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = database.DB.Exec(query, title, author, category, description, handler.Filename, coverHandler.Filename)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to save book to database")
		return
	}

	http.Redirect(w, r, "/admin.html", http.StatusSeeOther)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, author, category, description, book_file_path FROM books")
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to fetch books")
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Description, &book.BookFilePath)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Error scanning books")
			return
		}
		books = append(books, book)
	}

	utils.JSON(w, http.StatusOK, books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var book models.Book
	err = database.DB.QueryRow("SELECT id, title, author, category, description, book_file_path FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Category, &book.Description, &book.BookFilePath)

	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusNotFound, "Book not found")
		return
	} else if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.JSON(w, http.StatusOK, book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var filename string
	err = database.DB.QueryRow("SELECT book_file_path FROM books WHERE id = $1", id).Scan(&filename)
	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusNotFound, "Book not found")
		return
	}

	filePath := filepath.Join("../book", filename)
	_ = os.Remove(filePath)

	_, err = database.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to delete book from database")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Book deleted successfully"})
}
