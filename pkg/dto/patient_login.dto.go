package dto

type PatientLoginRequestDto struct {
	HospitalID string `json:"hospital_id" required:"required"`
	Password string `json:"password" required:"required"`
}

type PatientLoginResponseDto struct {
	AccessToken string `json:"access_token"`
}