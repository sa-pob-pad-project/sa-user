package service

import (
	"context"
	"fmt"
	"user-service/pkg/apperr"
	"user-service/pkg/clients"
	"user-service/pkg/constants"
	contextUtils "user-service/pkg/context"
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
	patientRepository *repository.PatientRepository
	doctorRepository  *repository.DoctorRepository
	userClient        *clients.UserClient
	jwtService        *jwt.JwtService
}

func NewUserService(
	db *gorm.DB,
	userRepo *repository.UserRepository,
	patientRepo *repository.PatientRepository,
	doctorRepo *repository.DoctorRepository,
	userClient *clients.UserClient,
	jwtService *jwt.JwtService,
) *UserService {
	return &UserService{
		db:                db,
		userRepository:    userRepo,
		patientRepository: patientRepo,
		doctorRepository:  doctorRepo,
		userClient:        userClient,
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
		return &dto.PatientRegisterResponseDto{}, apperr.New(apperr.CodeInternal, "hash password failed", err)
	}
	user.Password = hashedPassword

	tx := s.db.Begin()
	if tx.Error != nil {
		return &dto.PatientRegisterResponseDto{}, apperr.New(apperr.CodeInternal, "begin transaction failed", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userRepo := repository.NewUserRepository(tx)
	patientRepo := repository.NewPatientRepository(tx)

	if err := userRepo.Create(ctx, user); err != nil {
		tx.Rollback()
		return &dto.PatientRegisterResponseDto{}, apperr.New(apperr.CodeInternal, "create user failed", err)
	}

	if err := patientRepo.Create(ctx, patient); err != nil {
		tx.Rollback()
		return &dto.PatientRegisterResponseDto{}, apperr.New(apperr.CodeInternal, "create patient failed", err)
	}

	if err := tx.Commit().Error; err != nil {
		return &dto.PatientRegisterResponseDto{}, apperr.New(apperr.CodeInternal, "commit transaction failed", err)
	}

	return &dto.PatientRegisterResponseDto{Message: "User registered successfully"}, nil
}

func (s *UserService) PatientLogin(ctx context.Context, body *dto.PatientLoginRequestDto) (*dto.PatientLoginResponseDto, error) {
	patient, err := s.patientRepository.FindByHospitalID(ctx, body.HospitalID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to find patient", err)
	}
	if patient == nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid credentials", nil)
	}

	user, err := s.userRepository.FindByID(ctx, patient.UserID.String())
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to find user", err)
	}
	if user == nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid credentials", nil)
	}

	ok, err := utils.VerifyPassword(body.Password, user.Password)
	if !ok || err != nil {
		fmt.Println(ok, err)
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid credentials", err)
	}
	// sign token
	token, err := s.jwtService.GenerateToken(user.ID.String(), constants.RolePatient)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to generate token", err)
	}
	return &dto.PatientLoginResponseDto{
		AccessToken: token,
	}, nil
}

func (s *UserService) DoctorLogin(ctx context.Context, body *dto.DoctorLoginRequestDto) (*dto.DoctorLoginResponseDto, error) {
	doctor, err := s.doctorRepository.FindByUsername(ctx, body.Username)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to find doctor", err)
	}
	fmt.Println("doctor", doctor)
	if doctor == nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid credentials", nil)
	}

	user, err := s.userRepository.FindByID(ctx, doctor.UserID.String())
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to find user", err)
	}
	if user == nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid credentials", nil)
	}

	ok, err := utils.VerifyPassword(body.Password, user.Password)
	if !ok || err != nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid credentials", err)
	}
	// sign token
	token, err := s.jwtService.GenerateToken(user.ID.String(), constants.RoleDoctor)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to generate token", err)
	}
	return &dto.DoctorLoginResponseDto{
		AccessToken: token,
	}, nil

}

func (s *UserService) GetProfileByID(ctx context.Context) (*dto.GetProfileResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "user not found", err)
	}

	patient, err := s.patientRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "patient profile not found", err)
	}

	res := &dto.GetProfileResponseDto{
		ID:               user.ID.String(),
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Gender:           user.Gender,
		PhoneNumber:      user.PhoneNumber,
		HospitalID:       patient.HospitalID,
		BirthDate:        patient.BirthDate,
		IDCardNumber:     patient.IDCardNumber,
		Address:          patient.Address,
		Allergies:        patient.Allergies,
		EmergencyContact: patient.EmergencyContact,
		BloodType:        patient.BloodType,
	}
	return res, nil
}

func (s *UserService) GetPatientByID(ctx context.Context, patientID string) (*dto.GetProfileResponseDto, error) {
	user, err := s.userRepository.FindByID(ctx, patientID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "user not found", err)
	}

	patient, err := s.patientRepository.FindByUserID(ctx, patientID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "patient not found", err)
	}
	if patient == nil {
		return nil, apperr.New(apperr.CodeNotFound, "patient data not found", nil)
	}

	res := &dto.GetProfileResponseDto{
		ID:               user.ID.String(),
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Gender:           user.Gender,
		PhoneNumber:      user.PhoneNumber,
		HospitalID:       patient.HospitalID,
		BirthDate:        patient.BirthDate,
		IDCardNumber:     patient.IDCardNumber,
		Address:          patient.Address,
		Allergies:        patient.Allergies,
		EmergencyContact: patient.EmergencyContact,
		BloodType:        patient.BloodType,
	}
	return res, nil
}

func (s *UserService) GetDoctorByID(ctx context.Context, doctorID string) (*dto.GetDoctorProfileResponseDto, error) {
	user, err := s.userRepository.FindByID(ctx, doctorID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "user not found", err)
	}

	doctor, err := s.doctorRepository.FindByUserID(ctx, doctorID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "doctor not found", err)
	}
	if doctor == nil {
		return nil, apperr.New(apperr.CodeNotFound, "doctor data not found", nil)
	}

	res := &dto.GetDoctorProfileResponseDto{
		ID:              user.ID.String(),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Gender:          user.Gender,
		PhoneNumber:     user.PhoneNumber,
		Username:        doctor.Username,
		Specialty:       doctor.Specialty,
		Bio:             doctor.Bio,
		YearsExperience: doctor.YearsExperience,
	}
	return res, nil
}

func (s *UserService) UpdateProfileByID(ctx context.Context, body *dto.UpdatePatientProfileRequestDto) (*dto.UpdatePatientProfileResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, apperr.New(apperr.CodeInternal, "begin transaction failed", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userRepo := repository.NewUserRepository(tx)
	patientRepo := repository.NewPatientRepository(tx)

	user, err := userRepo.FindByID(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, apperr.New(apperr.CodeNotFound, "user not found", err)
	}
	if err != nil {
		tx.Rollback()
		return nil, apperr.New(apperr.CodeInternal, "failed to find user", err)
	}

	patient, err := patientRepo.FindByUserID(ctx, userID)
	if err != nil {
		tx.Rollback()
		return nil, apperr.New(apperr.CodeInternal, "failed to find patient", err)
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
		patient.BirthDate = body.BirthDate
	}
	if body.IDCardNumber != nil {
		patient.IDCardNumber = body.IDCardNumber
	}
	if body.Address != nil {
		patient.Address = body.Address
	}
	if body.Allergies != nil {
		patient.Allergies = body.Allergies
	}
	if body.EmergencyContact != nil {
		patient.EmergencyContact = body.EmergencyContact
	}
	if body.BloodType != nil {
		patient.BloodType = body.BloodType
	}

	// Save user and patient
	if err := userRepo.Update(ctx, user); err != nil {
		tx.Rollback()
		return nil, apperr.New(apperr.CodeInternal, "update user failed", err)
	}
	if err := patientRepo.Update(ctx, patient); err != nil {
		tx.Rollback()
		return nil, apperr.New(apperr.CodeInternal, "update patient failed", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, apperr.New(apperr.CodeInternal, "commit transaction failed", err)
	}

	return &dto.UpdatePatientProfileResponseDto{Message: "Profile updated successfully"}, nil
}

func (s *UserService) GetDoctorsByIDs(ctx context.Context, doctorIDs []string) ([]*dto.GetDoctorProfileResponseDto, error) {
	if len(doctorIDs) == 0 {
		return []*dto.GetDoctorProfileResponseDto{}, nil
	}
	doctors, err := s.userRepository.FindManyDoctorsByIDs(ctx, doctorIDs)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to find doctors", err)
	}
	var result []*dto.GetDoctorProfileResponseDto
	for _, user := range doctors {
		if user.Doctor == nil {
			continue
		}
		doctorDto := &dto.GetDoctorProfileResponseDto{
			ID:              user.ID.String(),
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Gender:          user.Gender,
			PhoneNumber:     user.PhoneNumber,
			Username:        user.Doctor.Username,
			Specialty:       user.Doctor.Specialty,
			Bio:             user.Doctor.Bio,
			YearsExperience: user.Doctor.YearsExperience,
		}
		result = append(result, doctorDto)
	}
	return result, nil
}

func (s *UserService) GetPatientsByIDs(ctx context.Context, patientIDs []string) ([]*dto.GetProfileResponseDto, error) {
	if len(patientIDs) == 0 {
		return []*dto.GetProfileResponseDto{}, nil
	}
	patients, err := s.userRepository.FindManyPatientsByIDs(ctx, patientIDs)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to find patients", err)
	}
	var result []*dto.GetProfileResponseDto
	for _, user := range patients {
		if user.Patient == nil {
			continue
		}
		patientDto := &dto.GetProfileResponseDto{
			ID:               user.ID.String(),
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
		result = append(result, patientDto)
	}
	return result, nil
}

func (s *UserService) GetAllDoctors(ctx context.Context) ([]*dto.GetDoctorProfileResponseDto, error) {
	doctors, err := s.doctorRepository.FindAll(ctx)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to find doctors", err)
	}
	var result []*dto.GetDoctorProfileResponseDto
	for _, doctor := range doctors {
		doctorDto := &dto.GetDoctorProfileResponseDto{
			ID:              doctor.User.ID.String(),
			FirstName:       doctor.User.FirstName,
			LastName:        doctor.User.LastName,
			Gender:          doctor.User.Gender,
			PhoneNumber:     doctor.User.PhoneNumber,
			Username:        doctor.Username,
			Specialty:       doctor.Specialty,
			Bio:             doctor.Bio,
			YearsExperience: doctor.YearsExperience,
		}
		result = append(result, doctorDto)
	}
	return result, nil
}
