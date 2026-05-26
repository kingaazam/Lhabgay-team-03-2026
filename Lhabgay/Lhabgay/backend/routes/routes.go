package routes

import (
	"net/http"

	"backend/controllers"
	"backend/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", controllers.Login).Methods(http.MethodPost)
	router.HandleFunc("/user/signup", controllers.Signup).Methods(http.MethodPost)
	router.HandleFunc("/logout", controllers.Logout).Methods(http.MethodGet)

	router.Handle("/admin/books/upload", middleware.RequireAdmin(http.HandlerFunc(controllers.UploadBook))).Methods(http.MethodPost)
	router.Handle("/admin/books/{id:[0-9]+}", middleware.RequireAdmin(http.HandlerFunc(controllers.DeleteBook))).Methods(http.MethodDelete)

	router.HandleFunc("/books", controllers.GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/books/{id:[0-9]+}", controllers.GetBook).Methods(http.MethodGet)

	router.HandleFunc("/files/{filepath:.+}", controllers.ServeFile).Methods(http.MethodGet)
}
