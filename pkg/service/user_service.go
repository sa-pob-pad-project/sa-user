package service

import (
	"context"
	"errors"
	"fmt"
	"user-service/pkg/apperr"
	"user-service/pkg/constants"
	"user-service/pkg/dto"
	"user-service/pkg/jwt"
	"user-service/pkg/models"
	"user-service/pkg/repository"
	"user-service/pkg/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db             *gorm.DB
	userRepository *repository.UserRepository
	jwtService     *jwt.JwtService
}

func NewUserService(db *gorm.DB, userRepo *repository.UserRepository, jwtService *jwt.JwtService) *UserService {
	return &UserService{
		db:             db,
		userRepository: userRepo,
		jwtService:     jwtService,
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
		return &dto.PatientRegisterResponseDto{}, apperr.New(apperr.CodeInternal, "hash password failed", err)
	}
	user.Password = hashedPassword

	_, err = s.userRepository.Transaction(ctx, func(repo *repository.UserRepository) (interface{}, error) {
		if err := repo.CreateUser(ctx, user); err != nil {
			return nil, apperr.New(apperr.CodeInternal, "create user failed", err)
		}
		if err := repo.CreatePatient(ctx, patient); err != nil {
			return nil, apperr.New(apperr.CodeInternal, "create patient failed", err)
		}
		return nil, nil
	})

	if err != nil {
		return &dto.PatientRegisterResponseDto{}, apperr.New(apperr.CodeInternal, "transaction failed", err)
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
	if err != nil {
		return nil, err
	}
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
	return res, nil
}

func (s *UserService) UpdateProfileByID(ctx context.Context, userID string, body *dto.UpdatePatientProfileRequestDto) (*dto.UpdatePatientProfileResponseDto, error) {
	result, err := s.userRepository.Transaction(ctx, func(repo *repository.UserRepository) (interface{}, error) {
		user, err := repo.FindPatientByID(ctx, userID)
		if err != nil {
			return nil, apperr.New(apperr.CodeBadRequest, "user not found", nil)
		}

		// Update user fields if provided
		if body.FirstName != nil {
			user.FirstName = *body.FirstName
		}
		if body.LastName != nil {
			user.LastName = *body.LastName
		}
		if body.PhoneNumber != nil {
			user.PhoneNumber = *body.PhoneNumber
		}

		// Update patient fields if provided
		if body.BirthDate != nil {
			user.Patient.BirthDate = body.BirthDate
		}
		if body.IDCardNumber != nil {
			user.Patient.IDCardNumber = body.IDCardNumber
		}
		if body.Address != nil {
			user.Patient.Address = body.Address
		}
		if body.Allergies != nil {
			user.Patient.Allergies = body.Allergies
		}
		if body.EmergencyContact != nil {
			user.Patient.EmergencyContact = body.EmergencyContact
		}
		if body.BloodType != nil {
			user.Patient.BloodType = body.BloodType
		}

		// Save user and patient
		if err := repo.UpdateUser(ctx, user); err != nil {
			return nil, apperr.New(apperr.CodeInternal, "update user failed", err)
		}
		if err := repo.UpdatePatient(ctx, user.Patient); err != nil {
			return nil, apperr.New(apperr.CodeInternal, "update patient failed", err)
		}

		return &dto.UpdatePatientProfileResponseDto{Message: "Profile updated successfully"}, nil
	})

	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "transaction failed", err)
	}

	return result.(*dto.UpdatePatientProfileResponseDto), nil
}
