package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexsasharegan/go-rss-agg-example/internal/database"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("Missing environment variable: PORT")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("Missing environment variable: DB_URL")
	}

	conn, err := sql.Open("postgres", dbURL+"?sslmode=disable")
	if err != nil {
		log.Fatalln("Cannot connect to database", dbURL)
	}

	api := apiConfig{
		DB: database.New(conn),
	}

	startScraping(api.DB, 16, 1*time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	router.Mount("/api/v1", v1Router)

	v1Router.Get("/healthz", healthzHandler)
	v1Router.Get("/errorz", errorzHandler)
	v1Router.Post("/users", api.handlerCreateUser)
	v1Router.Get("/feeds", api.handlerGetFeeds)
	// Authenticated Routes ------------------------------------------------
	v1Router.Get("/posts", api.middlewareAuth(api.handlerGetPostsByUser))
	v1Router.Get("/users", api.middlewareAuth(api.handlerGetUser))
	v1Router.Post("/feeds", api.middlewareAuth(api.handlerCreateFeed))
	v1Router.Post("/follows", api.middlewareAuth(api.handlerCreateFeedFollow))
	v1Router.Get("/follows", api.middlewareAuth(api.handlerGetFeedFollows))
	v1Router.Delete("/follows/{feedFollowID}", api.middlewareAuth(api.handlerDeleteFeedFollow))

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", port),
	}

	log.Printf("Server running on :%s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
