package repository

import (
	"context"
	"user-service/pkg/models"

	"gorm.io/gorm"
)

type DoctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) *DoctorRepository {
	return &DoctorRepository{
		db: db,
	}
}

func (r *DoctorRepository) Transaction(ctx context.Context, fn func(repo *DoctorRepository) (interface{}, error)) (interface{}, error) {
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

func (r *DoctorRepository) withTx(tx *gorm.DB) *DoctorRepository {
	return &DoctorRepository{db: tx}
}

func (r *DoctorRepository) FindByUserID(ctx context.Context, userID string) (*models.Doctor, error) {
	var doctor models.Doctor
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&doctor).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *DoctorRepository) FindByUsername(ctx context.Context, username string) (*models.Doctor, error) {
	var doctor models.Doctor
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&doctor).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *DoctorRepository) Create(ctx context.Context, doctor *models.Doctor) error {
	if err := r.db.WithContext(ctx).Create(doctor).Error; err != nil {
		return err
	}
	return nil
}

func (r *DoctorRepository) Update(ctx context.Context, doctor *models.Doctor) error {
	if err := r.db.WithContext(ctx).Save(doctor).Error; err != nil {
		return err
	}
	return nil
}

func (r *DoctorRepository) FindManyByIDs(ctx context.Context, doctorIDs []string) ([]*models.Doctor, error) {
	var doctors []*models.Doctor
	if err := r.db.WithContext(ctx).Where("user_id IN ?", doctorIDs).Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}
