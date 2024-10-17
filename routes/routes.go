package routes

import (
	"rr/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, Handler *handler.BannerHandler) {
	Admin := app.Group("api/admin/")
	Admin.Static("uploads", "./uploads")
	Admin.Post("banners", Handler.Create)
	Admin.Get("banners", Handler.GetPaginated)
	Admin.Get("banners/:id", Handler.GetByID)
	Admin.Delete("banners/:id", Handler.Delete)
	Admin.Put("banners/:id", Handler.Update)

}

func SetupEmployerRoutes(app *fiber.App, Handler *handler.EmployerHandler) {
	Employer := app.Group("api/admin/")
	Employer.Static("uploads", "./uploads")
	Employer.Post("employers", Handler.Create)
	Employer.Get("employers/:id", Handler.GetByID)
	Employer.Get("employers/", Handler.GetPaginated)
	Employer.Delete("employers/:id", Handler.Delete)
	Employer.Put("employers/:id", Handler.Update)

}
func SetupNewsRoutes(app *fiber.App, Handler *handler.NewsHandler) {
	News := app.Group("api/admin/")
	News.Static("uploads", "./uploads")
	News.Post("news", Handler.Create)
	News.Get("news/:id", Handler.GetByID)
	News.Get("news/", Handler.GetPaginated)
	News.Delete("news/:id", Handler.Delete)
	News.Put("news/:id", Handler.Update)

}
func SetupMediaRoutes(app *fiber.App, Handler *handler.MediaHandler) {
	News := app.Group("api/admin/")
	News.Static("uploads", "./uploads")
	News.Post("media", Handler.Create)
	News.Get("media/:id", Handler.GetByID)
	News.Get("media/", Handler.GetPaginated)
	News.Delete("media/:id", Handler.Delete)
	News.Put("media/:id", Handler.Update)

}
func AdminRoutes(app *fiber.App, adminHandler *handler.AdminHandler) {
	app.Post("/register", adminHandler.Register)
	app.Post("/login", adminHandler.Login)
}
