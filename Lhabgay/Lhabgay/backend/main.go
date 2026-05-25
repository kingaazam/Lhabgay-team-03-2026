package main

import (
	"log"
	"net/http"
	"os" // <-- 1. Added this so we can read environment variables

	"lhabgay/backend/database"
	"lhabgay/backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Login.html")
	}).Methods(http.MethodGet)

	// Frontend HTML files and image/book folders live in the project root.
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	// 2. Dynamic Port Handling for Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback to 8080 on your local machine
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router)) // <-- 3. Uses Render's port
}
