package main

import (
	"rr/database"
	"rr/handler"
	"rr/middleware"
	repository "rr/repostory"
	"rr/routes"
	"rr/service"
	"rr/setup"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupApp(db *gorm.DB, secret string) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})

	// Admin handler-leri belläp goýmak
	adminRepo := repository.NewAdminGormRepo(db)
	adminService := service.NewAdminService(adminRepo, secret)
	adminHandler := handler.NewAdminHandler(adminService)

	routes.AdminRoutes(app, adminHandler)

	// Middleware JWT barlagy üçin
	app.Use("/admin", middleware.JWTMiddleware(secret))

	return app
}

func main() {
	database.ConnectDB()
	app := SetupApp(database.DB, "your_secret")

	// Handler-leri ulanyp, routlary sazlaýar
	Handler := setup.SetupServices(database.DB)
	HandlerEmployer := setup.SetupEmployerServices(database.DB)
	HandlerNews := setup.SetupNewsServices(database.DB)
	HandlerMedia := setup.SetupMediaServices(database.DB)

	// Routlary kesgitlemek
	routes.SetupRoutes(app, Handler)
	routes.SetupEmployerRoutes(app, HandlerEmployer)
	routes.SetupNewsRoutes(app, HandlerNews)
	routes.SetupMediaRoutes(app, HandlerMedia)

	// Serweri diňlemäge başlaýar
	app.Listen(":5000")
}
