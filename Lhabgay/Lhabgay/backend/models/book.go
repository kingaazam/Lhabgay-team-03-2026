package models

type Book struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	BookFilePath string `json:"book_file_path"`
	// ADD THIS LINE RIGHT HERE:
	CoverImage string `json:"cover_image"`
}
