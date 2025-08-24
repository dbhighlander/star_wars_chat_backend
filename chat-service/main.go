package main

import (
	"log"
	"net/http"
	"os"
	"star_wars/m/internal/appMiddleware"
	"star_wars/m/internal/db"

	"star_wars/m/internal/routes"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {

	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MySQL successfully!")

	// Create chi router
	r := chi.NewRouter()

	// Set up CORS middleware
	frontendURL := os.Getenv("FRONTEND_URL")

	//add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(appMiddleware.ApiAuth)

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
