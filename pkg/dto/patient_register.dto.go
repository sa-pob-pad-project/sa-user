package dto

import "time"

type PatientRegisterPatientRequestDto struct {
	Password    string `json:"password" validate:"required,min=6"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Gender      string `json:"gender" validate:"required,oneof='male' 'female' 'other'"`
	PhoneNumber string `json:"phone_number"`
	// Patient specific fields
	HospitalID       string     `json:"hospital_id" validate:"required"`
	BirthDate        *time.Time `json:"birth_date,omitempty"`
	IDCardNumber     *string    `json:"id_card_number,omitempty" validate:"omitempty,len=13,numeric"`
	Address          *string    `json:"address,omitempty"`
	Allergies        *string    `json:"allergies,omitempty"`
	EmergencyContact *string    `json:"emergency_contact,omitempty"`
	BloodType        *string    `json:"blood_type,omitempty"`
}

type PatientRegisterResponseDto struct {
	Message string `json:"message"`
}
