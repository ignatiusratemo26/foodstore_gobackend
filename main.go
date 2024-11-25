package main

import (
	"log"
	"net/http"

	"go_backend/data"
	"go_backend/routes"

	"github.com/rs/cors"
)

func main() {
	data.InitMongo()
	router := routes.SetupRouter()

	// Add user routes
	routes.UserRoutes(router)

	// Add food routes
	routes.SetupFoodsRouter(router)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "X-Requested-With", "Origin"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(router)

	log.Println("Server running on:5000")
	http.ListenAndServe(":4000", handler)
}
