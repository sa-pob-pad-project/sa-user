package dto

import "time"

type UpdatePatientProfileRequestDto struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	// Patient specific fields
	BirthDate        *time.Time `json:"birth_date,omitempty"`
	IDCardNumber     *string    `json:"id_card_number,omitempty" validate:"omitempty,len=13,numeric"`
	Address          *string    `json:"address,omitempty"`
	Allergies        *string    `json:"allergies,omitempty"`
	EmergencyContact *string    `json:"emergency_contact,omitempty"`
	BloodType        *string    `json:"blood_type,omitempty"`
}

type UpdatePatientProfileResponseDto struct {
	Message string `json:"message"`
}
