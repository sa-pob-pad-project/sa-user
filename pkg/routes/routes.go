package routes

import (
	// "user-service/pkg/context"
	// "user-service/pkg/dto"
	_ "user-service/docs"
	"user-service/pkg/handlers"
	"user-service/pkg/jwt"
	"user-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, jwtSvc *jwt.JwtService) {

	api := app.Group("/api")
	user := api.Group("/user")
	user.Get("/swagger/*", swagger.HandlerDefault)

	v1 := user.Group("/v1")

	v1.Post("/patient/register",
		userHandler.PatientRegister) // TODO add validation middleware
	v1.Post("/patient/login", userHandler.PatientLogin)
	v1.Post("/doctor/login", userHandler.DoctorLogin)

	v1.Use(middleware.JwtMiddleware(jwtSvc))
	v1.Get("/patient/me", userHandler.Profile)
	v1.Patch("/patient", userHandler.UpdatePatientProfile)

	v1.Get("/doctors", userHandler.GetAllDoctors)
	v1.Post("/doctors", userHandler.GetDoctorByIDs)
	v1.Post("/patients", userHandler.GetPatientByIDs)
}
