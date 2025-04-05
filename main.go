package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type URL struct {
	ID       uint   `gorm:"primaryKey"`
	LongURL  string `gorm:"not null"`
	ShortURL string `gorm:"uniqueIndex"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
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
	log.Println("received request to /create")
	longURL := r.URL.Query().Get("longURL") //extracts the query parameters from request URL
	if longURL == "" {
		log.Println("missing longURL  parameter")
		//send error message as bad request if URL if request URL is empty
		http.Error(w, "LongURL is required", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	log.Printf("generated short URL %s\n", shortURL)

	url := URL{LongURL: longURL, ShortURL: shortURL}
	result := db.Create(&url)
	log.Printf("created result...%v\n", result)

	if result.Error != nil {
		log.Printf("database error %v\n", result.Error)
		//send internal server error if ther's an error creating URL
		http.Error(w, "Internal Server error, failed to create URL", http.StatusInternalServerError)
	}

	log.Printf("successfully created short URL %s\n", shortURL)
	fmt.Fprintf(w, "Short URL: %s\n", shortURL)
}

func redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:] //extracts the path from request URL
	var url URL
	result := db.First(&url, "short_url = ?", shortURL)
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url.LongURL, http.StatusFound)
}

func main() {
	//setup the handlers
	http.HandleFunc("/create", createShortURL)
	http.HandleFunc("/", redirectToLongURL)

	//handle favicon requests
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	//start the server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	//check the result
}
