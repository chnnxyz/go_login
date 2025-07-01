package handlers

import (
	"context"
	"cyberia_auth/config"
	"cyberia_auth/helpers"
	"cyberia_auth/models"
	"cyberia_auth/proto/notificationpb"
	"cyberia_auth/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type Credentials struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
	RoleName *string `json:"roleName"`
}

func sendVerificationEmail(userID, email, username, token string) {
	conn, err := grpc.Dial(
		"localhost:50054", grpc.WithInsecure(),
		grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Printf("Could not connect to email service: %v", err)
		return
	}
	defer conn.Close()

	client := notificationpb.NewEmailServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.SendVerificationEmail(ctx, &notificationpb.VerificationRequest{
		UserId:   userID,
		Email:    email,
		Username: username,
		Token:    token,
	})
	if err != nil {
		log.Printf("Failed to send verification email: %v", err)
	}
}

func RegisterPromoter(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Default to "Promoter" if RoleName is nil or empty
	roleName := "Promoter"

	var role models.Role
	if err := config.DB.First(&role, "name = ?", roleName).Error; err != nil {
		http.Error(w, "Role not found", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password encryption failed", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: creds.Username,
		Email:    *creds.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Username taken or database error", http.StatusBadRequest)
		return
	}
	token, err := helpers.GenerateToken(user.Email, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	sendVerificationEmail(user.ID.String(), user.Email, user.Username, token)
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
	token, err := utils.GenerateJWT(
		user.ID.String(), user.Username, user.Role.Name, isSuper,
	)

	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// Send token
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
