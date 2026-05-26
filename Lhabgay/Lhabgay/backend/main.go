package main

import (
	"log"
	"net/http"
	"os"

	"backend/database"
	"backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Explicitly serve Login.html from the parent project root folder
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("../Login.html"); err == nil {
			http.ServeFile(w, r, "../Login.html")
			return
		}
		// Fallback container context path
		http.ServeFile(w, r, "Login.html")
	}).Methods(http.MethodGet)

	// Serve static UI assets (.html, .css, .js) from the parent directory root
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
