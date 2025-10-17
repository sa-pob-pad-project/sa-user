package dto

type GetDoctorsByIDsRequestDto struct {
	DoctorIDs []string `json:"doctor_ids" validate:"required,dive,required,uuid"`
}