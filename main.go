package main

import (
  "os"
  "log"
  "net/http"
  
  "github.com/joho/godotenv"
  "github.com/gorilla/mux"

  "cyberia_auth/config"
  "cyberia_auth/models"
  "cyberia_auth/handlers"
  "golang.org/x/crypto/bcrypt"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found — using system env vars")
    }
}

func createSuperuser() {
  var count int64
  config.DB.Model(&models.User{}).Where("is_super = ?", true).Count(&count)
  if count == 0 {
    password, _ := bcrypt.GenerateFromPassword(
      []byte(os.Getenv("SU_PASSWORD")), bcrypt.DefaultCost,
    )
    username := os.Getenv("SUPERUSER_NAME")
    config.DB.Create(&models.User{
       Username: username,
       Password: string(password),
       IsSuper:  true,
    })
    log.Println("SuperUser created")
  } else {
    log.Println("Superuser exists")
  }
}

func main() {
    config.InitDB()
    createSuperuser()

    r := mux.NewRouter()
    r.HandleFunc("/auth/register", handlers.Register).Methods("POST")
    r.HandleFunc("/auth/login", handlers.Login).Methods("POST")

    log.Println("Auth service running on :8081")
    http.ListenAndServe(":8081", r)
}
