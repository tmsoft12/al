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

	routes.SetupRoutes(app, Handler)
	app.Listen(":5000")
}
