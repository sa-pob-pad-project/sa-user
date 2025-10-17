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

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindManyByIDs(ctx context.Context, userIDs []string) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindManyDoctorsByIDs(ctx context.Context, doctorIDs []string) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).
		Preload("Doctor").
		Where("id IN ?", doctorIDs).
		Where("role = ?", "doctor").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindManyPatientsByIDs(ctx context.Context, patientIDs []string) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).
		Preload("Patient").
		Where("id IN ?", patientIDs).
		Where("role = ?", "patient").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
