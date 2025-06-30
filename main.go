package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"cyberia_auth/config"
	"cyberia_auth/handlers"
	"cyberia_auth/utils"

	"github.com/rs/cors"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found — using system env vars")
	}
}

func main() {
	config.InitDb()
	utils.CreateRoles()
	utils.CreateSuperuser()

	r := mux.NewRouter()
	r.HandleFunc("/auth/register", handlers.RegisterPromoter).Methods("POST")
	r.HandleFunc("/auth/login", handlers.Login).Methods("POST")
	// Set up CORS options
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://admin.cyberiacollective.com",
			"http://localhost:5173"}, // Adjust based on your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := corsOptions.Handler(r)
	log.Println("Auth service running on :8081")
	http.ListenAndServe(":8081", handler)
}
