package main

import (
	"log"
	"net/http"
	"os"
	"star_wars/m/internal/db"
	"star_wars/m/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPass := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	// log.Printf("Connecting to DB at %s:%s with user %s, db %s, pw %s", dbHost, dbPort, dbUser, dbName, dbPass)
	// TODO: GORM connection here

	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MySQL successfully!")

	// Create chi router
	r := chi.NewRouter()

	// Set up CORS middleware
	frontendURL := os.Getenv("FRONTEND_URL")
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Register routes
	routes.RegisterRoutes(r)

	// Start server
	port := "8080"
	addr := "0.0.0.0:" + port
	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
