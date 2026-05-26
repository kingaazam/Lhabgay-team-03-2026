package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"backend/database"
	"backend/routes"

	"github.com/gorilla/mux"
)

// locateStaticDir searches up and down the layout to find where Login.html is running
func locateStaticDir() string {
	// 1. Check current folder execution context
	if _, err := os.Stat("Login.html"); err == nil {
		return "."
	}
	// 2. Check parent directory execution context
	if _, err := os.Stat("../Login.html"); err == nil {
		return ".."
	}
	// 3. Deep tree lookup fallback for Render container structures
	var foundPath string
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == "Login.html" {
			foundPath = filepath.Dir(path)
			return filepath.SkipDir
		}
		return nil
	})
	if foundPath != "" {
		return foundPath
	}
	return "." // Absolute fallback default
}

func main() {
	database.ConnectDB()

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	staticDir := locateStaticDir()
	log.Printf("[Static Asset Watchdog] Serving web assets from directory context: %s", staticDir)

	// Route root path explicitly to the located Login.html position
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetFile := filepath.Join(staticDir, "Login.html")
		http.ServeFile(w, r, targetFile)
	}).Methods(http.MethodGet)

	// Map all sub-resource routes (.css, .js, image files) dynamically
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
