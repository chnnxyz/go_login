package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"cyberia_auth/config"
	"cyberia_auth/handlers"
	"cyberia_auth/models"

	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found — using system env vars")
	}
}

func createSuperuser() {
	// Ensure SuperUser role exists
	var role models.Role
	if err := config.DB.FirstOrCreate(&role, models.Role{Name: "SuperUser"}).Error; err != nil {
		log.Fatalf("Failed to ensure SuperUser role exists: %v", err)
	}

	// Check if any user has the SuperUser role
	var count int64
	if err := config.DB.Model(&models.User{}).Where("role_id = ?", role.ID).Count(&count).Error; err != nil {
		log.Fatalf("Failed to check for existing superuser: %v", err)
	}

	if count == 0 {
		// Create new superuser
		password, _ := bcrypt.GenerateFromPassword(
			[]byte(os.Getenv("SU_PASSWORD")), bcrypt.DefaultCost,
		)
		username := os.Getenv("SUPERUSER_NAME")

		user := models.User{
			Username: username,
			Password: string(password),
			RoleID:   role.ID,
		}

		if err := config.DB.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create superuser: %v", err)
		}

		log.Println("SuperUser created")
	} else {
		log.Println("SuperUser already exists")
	}
}

func main() {
	config.InitDb()
	createSuperuser()

	r := mux.NewRouter()
	r.HandleFunc("/auth/register", handlers.Register).Methods("POST")
	r.HandleFunc("/auth/login", handlers.Login).Methods("POST")
	// Set up CORS options
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://admin.cyberiacollective.com"}, // Adjust based on your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := corsOptions.Handler(r)
	log.Println("Auth service running on :8081")
	http.ListenAndServe(":8081", handler)
}
