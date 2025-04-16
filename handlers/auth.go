
package handlers

import (
  "cyberia_auth/config"
  "cyberia_auth/models"
  "cyberia_auth/utils"
  "encoding/json"
  "net/http"
  "golang.org/x/crypto/bcrypt"
)

type Credentials struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

// Register handles user registration, parsing a request into
func Register(w http.ResponseWriter, r *http.Request) {
  var creds Credentials
  json.NewDecoder(r.Body).Decode(&creds)

  hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(creds.Password), bcrypt.DefaultCost,
  )
  user := models.User{
      Username: creds.Username, Password: string(hashedPassword),
  }

  // config.DB.Create returns nil on no error,
  // TODO: handle different error types
  if err := config.DB.Create(&user).Error; err != nil {
    http.Error(w, "Username taken", http.StatusBadRequest)
    return
  }

  w.WriteHeader(http.StatusCreated)
}

// Login handles login. using http requests and responeses.
func Login(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    json.NewDecoder(r.Body).Decode(&creds)

    var user models.User
    
    // Similarly, GORM returns an error if the user is not found
    if err := config.DB.Where("username = ?", creds.Username).First(&user).Error; err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    // 
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    token, _ := utils.GenerateJWT(user.Username)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
