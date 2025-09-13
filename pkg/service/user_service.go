package service

import (
	"context"
	"errors"
	"fmt"
	"user-service/pkg/constants"
	"user-service/pkg/dto"
	"user-service/pkg/jwt"
	"user-service/pkg/models"
	"user-service/pkg/repository"
	"user-service/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db                *gorm.DB
	userRepository    *repository.UserRepository
	jwtService        *jwt.JwtService
}

func NewUserService(db *gorm.DB, userRepo *repository.UserRepository, jwtService *jwt.JwtService) *UserService {
	return &UserService{
		db:                db,
		userRepository:    userRepo,
		jwtService:        jwtService,
	}
}

func (s *UserService) Register(ctx context.Context, body *dto.PatientRegisterPatientRequestDto) (*dto.PatientRegisterResponseDto, error) {
	user := &models.User{
		ID:          utils.GenerateUUIDv7(),
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Gender:      body.Gender,
		Role:        constants.RolePatient,
		PhoneNumber: body.PhoneNumber,
	}
	patient := &models.Patient{
		UserID:           user.ID,
		HospitalID:       body.HospitalID,
		BirthDate:        body.BirthDate,
		IDCardNumber:     body.IDCardNumber,
		Address:          body.Address,
		Allergies:        body.Allergies,
		EmergencyContact: body.EmergencyContact,
		BloodType:        body.BloodType,
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return &dto.PatientRegisterResponseDto{}, err
	}
	user.Password = hashedPassword

	err = s.userRepository.Transaction(ctx, func(repo *repository.UserRepository) error {
		if err := repo.CreateUser(ctx, user); err != nil {
			return err
		}
		if err := repo.CreatePatient(ctx, patient); err != nil {
			return err
		}
		return nil
	})
	
	if err != nil {
		return &dto.PatientRegisterResponseDto{}, err
	}

	return &dto.PatientRegisterResponseDto{Message: "User registered successfully"}, nil
}

func (s *UserService) PatientLogin(ctx context.Context, body *dto.PatientLoginRequestDto) (*dto.PatientLoginResponseDto, error) {

	user, err := s.userRepository.FindByHospitalID(ctx, body.HospitalID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	ok, err := utils.VerifyPassword(body.Password, user.Password)
	if !ok || err != nil {
		fmt.Println(ok, err)
		return nil, errors.New("invalid credentials")
	}
	// sign token
	token, err := s.jwtService.GenerateToken(user.ID.String(), constants.RolePatient)
	if err != nil {
		return nil, err
	}
	// set token in cookie
	return &dto.PatientLoginResponseDto{
		AccessToken: token,
	}, nil
}

func (s *UserService) GetProfileByID(ctx context.Context, userID string) (*dto.GetProfileResponseDto, error) {
	user, err := s.userRepository.FindPatientByID(ctx, userID)
	res := &dto.GetProfileResponseDto{
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Gender:           user.Gender,
		PhoneNumber:      user.PhoneNumber,
		HospitalID:       user.Patient.HospitalID,
		BirthDate:        user.Patient.BirthDate,
		IDCardNumber:     user.Patient.IDCardNumber,
		Address:          user.Patient.Address,
		Allergies:        user.Patient.Allergies,
		EmergencyContact: user.Patient.EmergencyContact,
		BloodType:        user.Patient.BloodType,
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
