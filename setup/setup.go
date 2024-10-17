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
