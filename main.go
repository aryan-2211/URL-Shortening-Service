package main

import (
	"math/rand"
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

func main() {

}
