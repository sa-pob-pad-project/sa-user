package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	UserID           uuid.UUID      `json:"user_id" gorm:"primaryKey;type:uuid"`
	HospitalID       string        `json:"hospital_id" gorm:"unique,not null"`
	BirthDate        *time.Time     `json:"birth_date,omitempty"`
	IDCardNumber     *string        `json:"id_card_number,omitempty" gorm:"size:13"`
	Address          *string        `json:"address,omitempty"`
	Allergies        *string        `json:"allergies,omitempty"`
	EmergencyContact *string        `json:"emergency_contact,omitempty"`
	BloodType        *string        `json:"blood_type,omitempty" gorm:"size:5"`
	CreatedAt        time.Time      `json:"created_at" gorm:"default:now()"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"default:now()"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	User                   User                    `json:"user" gorm:"foreignKey:UserID;references:ID"`
	HealthcareEntitlements []HealthcareEntitlement `json:"healthcare_entitlements,omitempty" gorm:"many2many:user_healthcare_entitlement;foreignKey:UserID;joinForeignKey:PatientID;References:Name;joinReferences:HealthcareEntitlement"`
}

