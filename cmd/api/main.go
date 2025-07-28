package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/King0625/SD.urlshortener/internal/db"
	"github.com/King0625/SD.urlshortener/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	dsn := os.Getenv("POSTGRES_DSN")

	if err := db.RunMigration(dsn); err != nil {
		log.Fatalf("run migration error: %v", err)
	}

	conn, err := db.InitPostgres(dsn)
	if err != nil {
		log.Fatalf("Cannot init postgres conn: %v", err)
	}
	defer conn.Close(context.Background())

	queries := db.New(conn)
	h := &handler.UrlHandler{
		Queries: queries,
		Conn:    conn,
	}

	r := chi.NewRouter()
	r.Post("/shorten", h.ShortenURL)
	r.Get("/{code}", h.Redirect)

	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", r)
}
