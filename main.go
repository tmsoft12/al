package main

import (
	"rr/database"
	"rr/routes"
	"rr/setup"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})
	Handler := setup.SetupServices(database.DB)
	HandlerEmployer := setup.SetupEmployerServices(database.DB)
	HandlerNews := setup.SetupNewsServices(database.DB)
	HandlerMedia := setup.SetupMediaServices(database.DB)
	routes.SetupRoutes(app, Handler)
	routes.SetupEmployerRoutes(app, HandlerEmployer)
	routes.SetupNewsRoutes(app, HandlerNews)
	routes.SetupMediaRoutes(app, HandlerMedia)
	app.Listen(":5000")
}
