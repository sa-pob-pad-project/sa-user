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

func (r *UserRepository) Transaction(ctx context.Context, fn func(repo *UserRepository) (interface{}, error)) (interface{}, error) {
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

func (r *UserRepository) withTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) FindByHospitalID(ctx context.Context, hospitalID string) (*models.User, error) {
	var user models.User
	if err := r.db.Joins("JOIN patients ON patients.user_id = users.id").Where("patients.hospital_id = ?", hospitalID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindPatientByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Patient").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err

	}
	return &user, nil

}

func (r *UserRepository) FindDoctorByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Doctor").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
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

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {

		return err
	}
	return nil
}

func (r *UserRepository) UpdatePatient(ctx context.Context, patient *models.Patient) error {
	if err := r.db.WithContext(ctx).Save(patient).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindManyDoctorsByIDs(ctx context.Context, doctorIDs []string) ([]*models.User, error) {
	var doctors []*models.User
	if err := r.db.WithContext(ctx).Preload("Doctor").Where("id IN ?", doctorIDs).Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}
