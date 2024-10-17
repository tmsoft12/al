package setup

import (
	"rr/handler"
	"rr/repostrory"
	"rr/service"

	"gorm.io/gorm"
)

func SetupServices(db *gorm.DB) *handler.BannerHandler {
	bannerRepo := &repostrory.BannerRepo{DB: db}
	bannerService := &service.BannerService{Repo: bannerRepo}
	bannerHandler := &handler.BannerHandler{Service: bannerService}

	return bannerHandler
}

func SetupEmployerServices(db *gorm.DB) *handler.EmployerHandler {
	employerRepo := &repostrory.EmployerRepo{DB: db}
	employerService := &service.EmployerService{Repo: employerRepo}
	employerHandler := &handler.EmployerHandler{Service: employerService}

	return employerHandler
}

func SetupNewsServices(db *gorm.DB) *handler.NewsHandler {
	newsRepo := &repostrory.NewsRepo{DB: db}
	newsService := &service.NewsService{Repo: newsRepo}
	newsHandler := &handler.NewsHandler{Service: newsService}

	return newsHandler
}
