package models

import "time"

// Book matches the books table fields returned to the frontend.
type Book struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	CoverImagePath string    `json:"cover_image_path"`
	BookFilePath   string    `json:"book_file_path"`
	CreatedAt      time.Time `json:"created_at"`
}
