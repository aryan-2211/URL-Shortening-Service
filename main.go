package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type URL struct {
	ID       uint   `gorm:"primaryKey"`
	LongURL  string `gorm:"not null"`
	ShortURL string `gorm:"uniqueIndex"`
}

var db *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&URL{})
}

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortURL)
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	longURL := r.URL.Query().Get("longURL") //extracts the query parameters from request URL
	if longURL == "" {
		//send error message as bad request if URL if request URL is empty
		http.Error(w, "LongURL is required", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()

	url := URL{LongURL: longURL, ShortURL: shortURL}
	result := db.Create(&url)

	if result.Error != nil {
		//send internal server error if ther's an error creating URL
		http.Error(w, "Internal Server error, failed to create URL", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Short URL: %s\n", shortURL)
}

func main() {

}
