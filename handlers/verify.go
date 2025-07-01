package handlers

import (
	"cyberia_auth/helpers"
	"cyberia_auth/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := strings.TrimPrefix(r.URL.Path, "/auth/verify/")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	// Use the utility function you already wrote
	email, err := helpers.ConfirmToken(tokenStr, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
		return
	}

	// Lookup the user by email
	user, err := utils.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// If already verified
	if user.IsVerified {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email already verified"})
		return
	}

	// Mark user as verified
	user.IsVerified = true
	if err := utils.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully"})
}
