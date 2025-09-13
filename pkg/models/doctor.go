package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Doctor struct {
	UserID          uuid.UUID      `json:"user_id" gorm:"primaryKey;type:uuid"`
	Username        string         `json:"username" gorm:"unique;not null"`
	Specialty       *string        `json:"specialty,omitempty"`
	Bio             *string        `json:"bio,omitempty"`
	YearsExperience *int           `json:"years_experience,omitempty" gorm:"check:years_experience IS NULL OR years_experience >= 0"`
	CreatedAt       time.Time      `json:"created_at" gorm:"default:now()"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"default:now()"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
