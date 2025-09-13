package repository

import (
	"context"
	"user-service/pkg/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Transaction(ctx context.Context, fn func(repo *UserRepository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	repoWithTx := r.withTx(tx)

	if err := fn(repoWithTx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *UserRepository) withTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) FindByHospitalID(ctx context.Context, hospitalID string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Joins("JOIN patients ON patients.user_id = users.id").Where("patients.hospital_id = ?", hospitalID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindPatientByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Patient").Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

	}
	return &user, nil

}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreatePatient(ctx context.Context, patient *models.Patient) error {
	if err := r.db.WithContext(ctx).Create(patient).Error; err != nil {
		return err
	}
	return nil
}

