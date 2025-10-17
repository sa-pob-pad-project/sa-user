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

// DoctorLogin godoc
// @Summary Login a doctor
// @Description Authenticate a doctor and return access token with cookie
// @Tags doctors
// @Accept  json
// @Produce  json
// @Param doctor body dto.DoctorLoginRequestDto true "Doctor login credentials"
// @Success 200 {object} dto.DoctorLoginResponseDto "Doctor logged in successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Router /api/user/v1/doctor/login [post]
func (h *UserHandler) DoctorLogin(c *fiber.Ctx) error {
	var body dto.DoctorLoginRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}
	ctx := contextUtils.GetContext(c)
	res, err := h.userService.DoctorLogin(ctx, &body)
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
	fmt.Println("Hello from updatepatient profile")
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}
	fmt.Println("geting context")
	ctx := contextUtils.GetContext(c)
	// print ctx for debugging
	fmt.Println("UpdatePatientProfile endpoint hit, ctx:", ctx)

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

// GetDoctorByIDs godoc
// @Summary Get doctors by IDs
// @Description Get multiple doctor profiles by their IDs
// @Tags doctors
// @Accept  json
// @Produce  json
// @Param body body dto.GetDoctorsByIDsRequestDto true "Doctor IDs"
// @Success 200 {object} []dto.GetDoctorProfileResponseDto "Doctor profiles retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "Doctors not found"
// @Failure 500 {object} response.ErrorResponse "Failed to get doctor profiles"
// @Router /api/user/v1/doctors [post]
func (h *UserHandler) GetDoctorByIDs(c *fiber.Ctx) error {
	var body dto.GetDoctorsByIDsRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	ctx := contextUtils.GetContext(c)
	doctors, err := h.userService.GetDoctorsByIDs(ctx, body.DoctorIDs)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, doctors)
}

// GetPatientByIDs godoc
// @Summary Get patients by IDs
// @Description Get multiple patient profiles by their IDs
// @Tags patients
// @Accept  json
// @Produce  json
// @Param body body dto.GetPatientsByIDsRequestDto true "Patient IDs"
// @Success 200 {object} []dto.GetProfileResponseDto "Patient profiles retrieved successfully"
// @Failure 404 {object} response.ErrorResponse "Patients not found"
// @Failure 500 {object} response.ErrorResponse "Failed to get patient profiles"
// @Router /api/user/v1/patients [post]
func (h *UserHandler) GetPatientByIDs(c *fiber.Ctx) error {
	var body dto.GetPatientsByIDsRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	ctx := contextUtils.GetContext(c)
	patients, err := h.userService.GetPatientsByIDs(ctx, body.PatientIDs)
	if err != nil {
		return apperr.WriteError(c, err)
	}
	return response.OK(c, patients)
}
