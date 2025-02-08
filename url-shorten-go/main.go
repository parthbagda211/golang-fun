package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq" 
)

// URLShortener struct now includes a PostgreSQL database connection
type URLShortener struct {
	db            *sql.DB
	shortURLBase  string
	mutex         sync.Mutex
}

// NewURLShortener creates a new URLShortener with a PostgreSQL database connection
func NewURLShortener(base, dbURL string) (*URLShortener, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	return &URLShortener{
		db:           db,
		shortURLBase: base,
	}, nil
}

// GenerateShortURL 
func (s *URLShortener) GenerateShortURL(originalURL string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate a random 6-character short URL
	shortURL := generateRandomString(6)

	// Store the mapping in the database
	_, err := s.db.Exec("INSERT INTO url_mappings (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING", shortURL, originalURL)
	if err != nil {
		log.Fatalf("Error inserting URL mapping: %v", err)
	}

	return shortURL
}

// GetOriginalURL
func (s *URLShortener) GetOriginalURL(shortURL string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var originalURL string
	// Query the database for the original URL based on the short URL
	err := s.db.QueryRow("SELECT original_url FROM url_mappings WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err == sql.ErrNoRows {
		// If no row is found, return false
		return "", false
	} else if err != nil {
		log.Printf("Error retrieving original URL: %v", err)
		return "", false
	}

	return originalURL, true
}

// generateRandomString 
func generateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}

// handleShortenURL 
func (s *URLShortener) handleShortenURL(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	shortURL := s.GenerateShortURL(originalURL)
	shortenedURL := fmt.Sprintf("%s/%s", s.shortURLBase, shortURL)
	fmt.Fprintf(w, "Shortened URL: %s\n", shortenedURL)
}

// handleRedirectURL 
func (s *URLShortener) handleRedirectURL(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, exists := s.GetOriginalURL(shortURL)
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	
	dbURL := "postgres://postgres:12345@localhost:5432/url_shortener?sslmode=disable"

	// Initialize the URL shortener with PostgreSQL database connection
	urlShortener, err := NewURLShortener("http://localhost:5000", dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize URL shortener: %v", err)
	}

	// Register the HTTP routes
	http.HandleFunc("/shorten", urlShortener.handleShortenURL) // Shorten a URL
	http.HandleFunc("/", urlShortener.handleRedirectURL)      // Redirect to original URL

	// Start the server
	log.Println("Starting server on :5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
