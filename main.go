package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	apiConf := apiConfig{
		DB: database.New(conn),
	}

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
	v1Router.Post("/users", apiConf.handlerCreateUser)
	v1Router.Get("/feeds", apiConf.handlerGetFeeds)
	// Authenticated Routes ------------------------------------------------
	v1Router.Get("/users", apiConf.middlewareAuth(apiConf.handlerGetUser))
	v1Router.Post("/feeds", apiConf.middlewareAuth(apiConf.handlerCreateFeed))
	v1Router.Post("/follows", apiConf.middlewareAuth(apiConf.handlerCreateFeedFollow))
	v1Router.Get("/follows", apiConf.middlewareAuth(apiConf.handlerGetFeedFollows))
	v1Router.Delete("/follows/{feedFollowID}", apiConf.middlewareAuth(apiConf.handlerDeleteFeedFollow))

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", port),
	}

	log.Printf("Server running on :%s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
