package main

import (
	"rr/database"
	"rr/routes"
	"rr/setup"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()
	app := fiber.New()
	Handler := setup.SetupServices(database.DB)
	HandlerEmployer := setup.SetupEmployerServices(database.DB)
	HandlerNews := setup.SetupNewsServices(database.DB)
	routes.SetupRoutes(app, Handler)
	routes.SetupEmployerRoutes(app, HandlerEmployer)
	routes.SetupNewsRoutes(app, HandlerNews)
	app.Listen(":5000")
}
