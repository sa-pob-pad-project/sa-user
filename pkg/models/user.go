package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	PatientRole Role = "patient"
	DoctorRole  Role = "doctor"
	AdminRole   Role = "admin"
)


const (
	Male   string = "male"
	Female string = "female"
)

type User struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Password    string         `json:"-" gorm:"not null"`
	FirstName   string         `json:"first_name" gorm:"not null"`
	LastName    string         `json:"last_name" gorm:"not null"`
	Gender      string         `json:"gender" gorm:"type:gender_enum;not null"`
	PhoneNumber string         `json:"phone_number" gorm:"not null"`
	Role        Role           `json:"role" gorm:"type:roles;not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"default:now()"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Patient *Patient `json:"patient,omitempty" gorm:"foreignKey:UserID"`
	Doctor  *Doctor  `json:"doctor,omitempty" gorm:"foreignKey:UserID"`
	Admin   *Admin   `json:"admin,omitempty" gorm:"foreignKey:UserID"`
}
