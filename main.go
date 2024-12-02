package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

// Function to generate a short URL using MD5 hash
func generateshortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8] // Return first 8 characters of the hash
}

// Function to create and store a shortened URL
func createURL(OriginalURL string) string {
	shortURL := generateshortURL(OriginalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  OriginalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

// Function to retrieve the original URL from the database
func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

// Default handler
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the URL Shortener!")
}

// Handler for shortening a URL
func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	// Decode the incoming JSON request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate the short URL
	shortURL := createURL(data.URL)

	// Prepare the response
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	fmt.Println("Starting URL shortener service...")

	// Define routes
	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", ShortURLHandler)

	// Start the server
	fmt.Println("Server running on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
