package dto

type DoctorLoginRequestDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}


type DoctorLoginResponseDto struct {
	AccessToken string `json:"access_token"`
}