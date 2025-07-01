package utils

import (
	"cyberia_auth/config"
	"cyberia_auth/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func CreateRoles() {
	roleNames := []string{"SuperUser", "Promoter",
		"Security", "Seller", "User"}
	for _, role := range roleNames {
		var roletype models.Role
		if err := config.DB.FirstOrCreate(&roletype, models.Role{Name: role}).Error; err != nil {
			log.Fatalf("Failed to create or find role: %v", err)
		}
	}
}

func CreateSuperuser() {
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
		email := os.Getenv("SUPERUSER_EMAIL")

		user := models.User{
			Username: username,
			Password: string(password),
			Email:    email,
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

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}
