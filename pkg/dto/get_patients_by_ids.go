package dto

type GetPatientsByIDsRequestDto struct {
	PatientIDs []string `json:"patient_ids" validate:"required,dive,required,uuid"`
}
