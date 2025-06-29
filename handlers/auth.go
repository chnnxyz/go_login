package handlers

import (
	"cyberia_auth/config"
	"cyberia_auth/models"
	"cyberia_auth/utils"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	RoleID   *string `json:"roleId"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var roleID *uuid.UUID = nil
	if creds.RoleID != nil && *creds.RoleID != "" {
		parsedID, err := uuid.Parse(*creds.RoleID)
		if err != nil {
			http.Error(w, "Invalid role ID format", http.StatusBadRequest)
			return
		}
		// Validate role exists
		var role models.Role
		if err := config.DB.First(&role, "id = ?", parsedID).Error; err != nil {
			http.Error(w, "Role not found", http.StatusBadRequest)
			return
		}
		roleID = &parsedID
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)

	user := models.User{
		Username: creds.Username,
		Password: string(hashedPassword),
		RoleID:   *roleID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Username taken", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles login. using http requests and responeses.
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User

	if err := config.DB.Preload("Role").Where("username = ?", creds.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Determine if user is superuser based on role name
	isSuper := user.Role.Name == "SuperUser"

	// Generate token
	token, err := utils.GenerateJWT(user.Username, isSuper)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// Send token
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
