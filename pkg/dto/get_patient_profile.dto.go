package dto

import "time"

type GetProfileResponseDto struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Gender           string     `json:"gender"`
	PhoneNumber      string     `json:"phone_number"`
	HospitalID       string     `json:"hospital_id"`
	BirthDate        *time.Time `json:"birth_date"`
	IDCardNumber     *string    `json:"id_card_number"`
	Address          *string    `json:"address"`
	Allergies        *string    `json:"allergies"`
	EmergencyContact *string    `json:"emergency_contact"`
	BloodType        *string    `json:"blood_type"`
}
