package repostrory

import (
	"rr/domain"

	"gorm.io/gorm"
)

type MediaRepo struct {
	DB *gorm.DB
}

func (r *MediaRepo) Create(t *domain.Media) error {
	return r.DB.Create(t).Error
}
func (r *MediaRepo) FindAll() ([]domain.Media, error) {
	var media []domain.Media
	err := r.DB.Find(&media).Error
	return media, err
}

func (r *MediaRepo) FindByID(id uint) (*domain.Media, error) {
	var media domain.Media
	err := r.DB.First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}
func (r *MediaRepo) Delete(id uint) error {
	return r.DB.Delete(&domain.Media{}, id).Error
}
func (r *MediaRepo) Update(id uint, Media *domain.Media) error {
	return r.DB.Save(Media).Error
}
