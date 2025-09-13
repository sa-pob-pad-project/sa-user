package models

import (
	"github.com/google/uuid"
)

type HealthcareEntitlement struct {
	Name string `json:"name" gorm:"primaryKey;column:healthcare_entitlement"`

	Patients []Patient `json:"patients,omitempty" gorm:"many2many:user_healthcare_entitlement;foreignKey:Name;joinForeignKey:HealthcareEntitlement;References:UserID;joinReferences:PatientID"`
}

type UserHealthcareEntitlement struct {
	PatientID                 uuid.UUID `json:"patient_id" gorm:"primaryKey;type:uuid"`
	HealthcareEntitlementName string    `json:"healthcare_entitlement" gorm:"primaryKey;column:healthcare_entitlement"`

	Patient                  Patient               `json:"patient" gorm:"foreignKey:PatientID;references:UserID"`
	HealthcareEntitlementRef HealthcareEntitlement `json:"healthcare_entitlement_detail" gorm:"foreignKey:HealthcareEntitlementName;references:Name"`
}
