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

func generateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

func createurl(originalURL string) string {
	shorturl := generateShortURL(originalURL)
	id := shorturl
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shorturl,
		CreationDate: time.Now(),
	}
	return shorturl
}

func geturl(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("url not found")
	}
	return url, nil
}

func shorturlhandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}

	shortURL := createurl(data.URL)
	// Get Render’s dynamically assigned domain
domain := os.Getenv("RENDER_EXTERNAL_URL")
if domain == "" {
	domain = "https://url-shortener-1.onrender.com" // Fallback to your Render domain
}

response := map[string]string{
	"short_url": fmt.Sprintf("%s/redirect/%s", domain, shortURL),
}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirecturlhandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := geturl(id)
	if err != nil {
		http.Error(w, "url not found", http.StatusNotFound)
		return // Don't forget to return after error
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API routes
	http.HandleFunc("/shorten", shorturlhandler)
	http.HandleFunc("/redirect/", redirecturlhandler)

	// Root route serves frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	fmt.Println("Server running on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
