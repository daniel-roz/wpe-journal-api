package main

import (
	"api/database"
	"api/handlers"
	"api/repositories"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
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

	// Register API routes.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /journals/{category}", journalHandler.GetJournalsByCategory)
	mux.HandleFunc("GET /journals/{category}/{id}", journalHandler.GetJournalByID)
	mux.HandleFunc("POST /journals", journalHandler.CreateJournal)
	mux.HandleFunc("PATCH /journals/{id}", journalHandler.UpdateJournal)
	mux.HandleFunc("DELETE /journals/{id}", journalHandler.DeleteJournal)

	handler := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Server running on %s", addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		isAdmin := origin == "http://localhost:63342" ||
			origin == "http://127.0.0.1:63342"

		isPublic := origin == os.Getenv("FRONTEND_PUBLIC_URL")

		if isAdmin || isPublic {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PATCH, DELETE, OPTIONS",
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Public origin is read-only.
		if isPublic && r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
