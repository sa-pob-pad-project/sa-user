package models

import (
	"github.com/google/uuid"
)

type Admin struct {
	UserID   uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	Username string    `json:"username" gorm:"unique;not null"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
