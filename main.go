package main

import (
	"bytes"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	// "user-service/cmd"
	"user-service/pkg/clients"
	"user-service/pkg/config"
	dbpkg "user-service/pkg/db"
	"user-service/pkg/handlers"
	"user-service/pkg/jwt"
	"user-service/pkg/repository"
	"user-service/pkg/routes"
	service "user-service/pkg/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pressly/goose/v3"
)

//go:embed pkg/db/migrations/*.sql
var migrationsFS embed.FS

// migrateUp applies all up migrations using goose with embedded FS.
// Set env MIGRATE_ON_START=true (default) to run on startup.
func migrateUp(sqlDB *sql.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}
	// directory path is relative to the root of the embedded FS
	if err := goose.Up(sqlDB, "pkg/db/migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

// @title User API
// @description This is a sample server for a user API.
// @version 1.0
// @host localhost:5000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	os.Setenv("TZ", "Asia/Bangkok")
	config.LoadConfig()

	gormDB := dbpkg.Open(dbpkg.Config{
		Host:     config.Get("DB_HOST", "localhost"),
		Port:     config.GetInt("DB_PORT", 5432),
		User:     config.Get("DB_USER", "user"),
		Password: config.Get("DB_PASSWORD", "password"),
		Dbname:   config.Get("DB_NAME", "userdb"),
		Sslmode:  config.Get("DB_SSLMODE", "disable"),
	})

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("cannot get *sql.DB from gorm: %v", err)
	}
	// cmd.InitCmd()

	// Run migrations on start if enabled
	if config.Get("MIGRATE_ON_START", "true") == "true" {
		if err := migrateUp(sqlDB); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	}

	userServiceUrl := config.Get("USER_SERVICE_URL", "http://localhost:8000")
	userClient := clients.New(userServiceUrl)
	userRepository := repository.NewUserRepository(gormDB)
	patientRepository := repository.NewPatientRepository(gormDB)
	doctorRepository := repository.NewDoctorRepository(gormDB)
	jwtService := jwt.NewJwtService(
		config.Get("JWT_SECRET", "secret"),
		config.GetInt("JWT_TTL", 3600),
	)
	userService := service.NewUserService(gormDB, userRepository, patientRepository, doctorRepository, userClient, jwtService)
	userHandler := handlers.NewUserHandler(userService)
	validate := validator.New()

	app := fiber.New(fiber.Config{
		JSONDecoder: func(b []byte, v any) error {
			dec := json.NewDecoder(bytes.NewReader(b))
			dec.DisallowUnknownFields()
			if err := dec.Decode(v); err != nil {
				return fmt.Errorf("decode: %w", err)
			}
			if err := dec.Decode(new(struct{})); err != io.EOF {
				return fmt.Errorf("decode: trailing data")
			}

			rv := reflect.ValueOf(v)
			for rv.Kind() == reflect.Pointer {
				rv = rv.Elem()
			}

			if rv.Kind() == reflect.Struct {
				if err := validate.Struct(v); err != nil {
					return err
				}
			}
			return nil
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	routes.SetupRoutes(app, userHandler, jwtService)

	port := config.Get("APP_PORT", "8000")
	fmt.Println("Server is running on port " + port)
	// listen on all interfaces for containerized envs; change back to "localhost:"+port if desired
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
