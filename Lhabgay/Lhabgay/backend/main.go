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

	// 1. Look for Login.html in the current project directory context
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Login.html") // <-- Removed the "../../"
	}).Methods(http.MethodGet)

	// 2. Serve static asset files from the current folder context
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(".")))) // <-- Changed to "."
	// 2. Safely serve all other static asset files (HTML, CSS, JS, Images) from the parent directory
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../../"))))

	// 3. Dynamic Port Handling for Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback to 8080 on your local machine
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
