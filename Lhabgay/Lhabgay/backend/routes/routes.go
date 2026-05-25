package routes

import (
	"net/http"

	"lhabgay/backend/controllers"
	"lhabgay/backend/middleware"

	"github.com/gorilla/mux"
)

// RegisterRoutes attaches all API endpoints to the router.
func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", controllers.Login).Methods(http.MethodPost)
	router.HandleFunc("/user/signup", controllers.Signup).Methods(http.MethodPost)
	router.HandleFunc("/logout", controllers.Logout).Methods(http.MethodGet)

	router.Handle("/admin/books/upload", middleware.RequireAdmin(http.HandlerFunc(controllers.UploadBook))).Methods(http.MethodPost)
	router.Handle("/admin/books/{id:[0-9]+}", middleware.RequireAdmin(http.HandlerFunc(controllers.DeleteBook))).Methods(http.MethodDelete)

	router.HandleFunc("/books", controllers.GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/books/{id:[0-9]+}", controllers.GetBook).Methods(http.MethodGet)

	// File serving endpoints - must come before catch-all PathPrefix
	router.HandleFunc("/files/{filepath:.+}", controllers.ServeFile).Methods(http.MethodGet)
}
