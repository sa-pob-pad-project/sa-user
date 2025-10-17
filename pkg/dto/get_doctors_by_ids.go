package dto

type GetDoctorsByIDsRequestDto struct {
	DoctorIDs []string `json:"doctor_ids" validate:"required,min=1,dive,required,uuid"`
}