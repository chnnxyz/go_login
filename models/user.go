package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Password     string    `gorm:"not null"`
	RoleID       uuid.UUID `gorm:"type:uuid"` // Foreign key
	Role         Role      `gorm:"foreignKey:RoleID"`
	RegisteredBy uuid.UUID
	IsVerified   bool `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BeforeCreate hook to set UUID before inserting a new record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
