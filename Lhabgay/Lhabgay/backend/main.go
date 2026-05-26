package main

import (
	"log"
	"net/http"
	"os"

	"Lhabgay/backend/database"
	"Lhabgay/backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Explicitly serve Login.html at the root URL
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try serving from current directory
		if _, err := os.Stat("Login.html"); err == nil {
			http.ServeFile(w, r, "Login.html")
			return
		}
		// Fallback check if it sits one level up from backend folder locally
		http.ServeFile(w, r, "../Login.html")
	}).Methods(http.MethodGet)

	// Serve all other frontend static files (.html, .css, .js, images)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
