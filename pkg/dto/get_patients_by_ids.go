package dto

type GetPatientsByIDsRequestDto struct {
	PatientIDs []string `json:"patient_ids" validate:"required,min=1,dive,required,uuid"`
}
