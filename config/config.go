package config

import (
    "log"
    "cyberia_auth/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    var err error
    DB, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect to database:", err)
    }

    DB.AutoMigrate(&models.User{})
}
