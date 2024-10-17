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
	Employer.Post("employers", Handler.Create)
	Employer.Get("employers/:id", Handler.GetByID)
	Employer.Get("employers/", Handler.GetPaginated)
	Employer.Delete("employers/:id", Handler.Delete)
	Employer.Put("employers/:id", Handler.Update)

}
