package handlers

import (
	"fmt"
	"user-service/pkg/apperr"
	contextUtils "user-service/pkg/context"
	"user-service/pkg/dto"
	response "user-service/pkg/response"
	service "user-service/pkg/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService}
}

// Handler functions
// PatientRegister godoc
// @Summary Register a new patient
// @Description Register a new patient in the system
// @Tags patients
// @Accept  json
// @Produce  json
// @Param patient body dto.PatientRegisterPatientRequestDto true "Patient registration data"
// @Success 201 {object} dto.PatientRegisterResponseDto "Patient registered successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to register user"
// @Router /api/user/v1/patient/register [post]
func (h *UserHandler) PatientRegister(c *fiber.Ctx) error {
	fmt.Println("Register endpoint hit")
	var body dto.PatientRegisterPatientRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.userService.Register(ctx, &body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return response.Created(c, res)
}

// PatientLogin godoc
// @Summary Login a patient
// @Description Authenticate a patient and return access token
// @Tags patients
// @Accept  json
// @Produce  json
// @Param patient body dto.PatientLoginRequestDto true "Patient login credentials"
// @Success 200 {object} dto.PatientLoginResponseDto "Patient logged in successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Router /api/user/v1/patient/login [post]
func (h *UserHandler) PatientLogin(c *fiber.Ctx) error {
	var body dto.PatientLoginRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.userService.PatientLogin(ctx, &body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	// set token in cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		HTTPOnly: true,
		SameSite: "None",
	})
	return response.OK(c, res)
}

// Profile godoc
// @Summary Get patient profile
// @Description Get the profile information of the authenticated patient
// @Tags patients
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} dto.GetProfileResponseDto "Profile retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 500 {object} response.ErrorResponse "Failed to get user profile"
// @Router /api/user/v1/patient/me [get]
func (h *UserHandler) Profile(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	user, err := h.userService.GetProfileByID(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, user)
}

// UpdatePatientProfile godoc
// @Summary Update patient profile
// @Description Update the profile information of the authenticated patient
// @Tags patients
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param patient body dto.UpdatePatientProfileRequestDto true "Patient profile update data"
// @Success 200 {object} dto.UpdatePatientProfileResponseDto "Profile updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or user not found"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 500 {object} response.ErrorResponse "Failed to update user profile"
// @Router /api/user/v1/patient/me [put]
func (h *UserHandler) UpdatePatientProfile(c *fiber.Ctx) error {
	var body dto.UpdatePatientProfileRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.userService.UpdateProfileByID(ctx, &body)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, res)
}

// GetPatientByID godoc
// @Summary Get patient by ID
// @Description Get patient profile information by patient ID
// @Tags patients
// @Accept  json
// @Produce  json
// @Param id path string true "Patient ID"
// @Success 200 {object} dto.GetProfileResponseDto "Patient profile retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "Patient not found"
// @Failure 500 {object} response.ErrorResponse "Failed to get patient profile"
// @Router /api/user/v1/patient/{id} [get]
func (h *UserHandler) GetPatientByID(c *fiber.Ctx) error {
	patientID := c.Params("id")
	fmt.Println("GetPatientByID endpoint hit, patientID:", patientID)

	ctx := contextUtils.GetContext(c)
	patient, err := h.userService.GetPatientByID(ctx, patientID)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, patient)
}

// GetDoctorByID godoc
// @Summary Get doctor by ID
// @Description Get doctor profile information by doctor ID
// @Tags doctors
// @Accept  json
// @Produce  json
// @Param id path string true "Doctor ID"
// @Success 200 {object} dto.GetDoctorProfileResponseDto "Doctor profile retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "Doctor not found"
// @Failure 500 {object} response.ErrorResponse "Failed to get doctor profile"
// @Router /api/user/v1/doctor/{id} [get]
func (h *UserHandler) GetDoctorByID(c *fiber.Ctx) error {
	doctorID := c.Params("id")
	fmt.Println("GetDoctorByID endpoint hit, doctorID:", doctorID)

	ctx := contextUtils.GetContext(c)
	doctor, err := h.userService.GetDoctorByID(ctx, doctorID)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, doctor)
}
