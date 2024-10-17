package repository

import (
	"rr/domain"

	"gorm.io/gorm"
)

type AdminRepository interface {
	Register(admin *domain.Admin) error
	Login(username string) (*domain.Admin, error)
}

type AdminGormRepo struct {
	DB *gorm.DB
}

func NewAdminGormRepo(db *gorm.DB) *AdminGormRepo {
	return &AdminGormRepo{
		DB: db,
	}
}

func (repo *AdminGormRepo) Register(admin *domain.Admin) error {
	return repo.DB.Create(admin).Error
}

func (repo *AdminGormRepo) Login(username string) (*domain.Admin, error) {
	var admin domain.Admin
	if err := repo.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}
