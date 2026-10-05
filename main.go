package main

import (
	"api/database"
	"api/handlers"
	"api/repositories"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := database.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	journalRepository := repositories.NewJournalRepository(db)
	journalHandler := handlers.NewJournalHandler(journalRepository)

	http.HandleFunc("GET /health", handlers.HealthHandler)
	http.HandleFunc("GET /journals/{category}", journalHandler.GetJournalsByCategory)
	http.HandleFunc("GET /journals/{category}/{id}", journalHandler.GetJournalByID)

	http.HandleFunc("POST /journals", journalHandler.CreateJournal)

	http.HandleFunc("PATCH /journals/{id}", journalHandler.UpdateJournal)

	http.HandleFunc("DELETE /journals/{id}", journalHandler.DeleteJournal)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
