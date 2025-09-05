package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/King0625/SD.urlshortener/internal/db"
	"github.com/King0625/SD.urlshortener/internal/handler"
	"github.com/King0625/SD.urlshortener/internal/repository"
	"github.com/King0625/SD.urlshortener/internal/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	urlRepository := repository.NewUrlRepository(conn)
	urlService := service.NewUrlService(urlRepository)
	urlHandler := handler.UrlHandler{Service: urlService}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.GetHead)
	r.Use(middleware.Heartbeat("/"))
	r.Use(middleware.RealIP)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.ThrottleBacklog(20, 100, time.Second*10))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Post("/shorten", urlHandler.ShortenURL)
	r.Get("/{code}", urlHandler.Redirect)
	r.Delete("/{code}", urlHandler.DeleteUrlByCode)
	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", r)
}
