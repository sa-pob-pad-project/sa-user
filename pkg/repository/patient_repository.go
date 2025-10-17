package repository

import (
	"context"
	"user-service/pkg/models"

	"gorm.io/gorm"
)

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{
		db: db,
	}
}

func (r *PatientRepository) Transaction(ctx context.Context, fn func(repo *PatientRepository) (interface{}, error)) (interface{}, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	repoWithTx := r.withTx(tx)

	result, err := fn(repoWithTx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PatientRepository) withTx(tx *gorm.DB) *PatientRepository {
	return &PatientRepository{db: tx}
}

func (r *PatientRepository) FindByUserID(ctx context.Context, userID string) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&patient).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) FindByHospitalID(ctx context.Context, hospitalID string) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.WithContext(ctx).Where("hospital_id = ?", hospitalID).First(&patient).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepository) Create(ctx context.Context, patient *models.Patient) error {
	if err := r.db.WithContext(ctx).Create(patient).Error; err != nil {
		return err
	}
	return nil
}

func (r *PatientRepository) Update(ctx context.Context, patient *models.Patient) error {
	if err := r.db.WithContext(ctx).Save(patient).Error; err != nil {
		return err
	}
	return nil
}

func (r *PatientRepository) FindManyByIDs(ctx context.Context, patientIDs []string) ([]*models.Patient, error) {
	var patients []*models.Patient
	if err := r.db.WithContext(ctx).Where("user_id IN ?", patientIDs).Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}
