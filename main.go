package main

import (
	"database/sql"
	"log"
	"net/http"
	"tender-scraper/config"
	"tender-scraper/database"
	"tender-scraper/scraper"
)

var db *sql.DB

func main() {

	log.Println("Starting scraper service...")
	performScrapingTask()
	log.Println("Completed scraper service...")
	
}

func initializeDatabase(cfg *config.Config) {
	var err error
	db, err = database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	log.Println("Successfully connected to the database. All good!")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust this to allow requests from specific domains
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func performScrapingTask() {		

	tenderDetails, err := scraper.ScrapeReviews()
	if err != nil {
		log.Fatal("Failed to scrape reviews: ", err)
	}
	print("tenderDetails", tenderDetails)

	log.Println("Scraping task completed.")
}


