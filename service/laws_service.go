package service

import (
	"rr/domain"
	repository "rr/repostory"
)

type LawsService struct {
	Repo *repository.LawsRepo
}

func (s *LawsService) Create(t *domain.Laws) error {

	return s.Repo.Create(t)
}
func (s *LawsService) GetAll() ([]domain.Laws, error) {
	return s.Repo.FindAll()
}
func (s *LawsService) GetPaginated(page int, limit int) ([]domain.Laws, int64, error) {
	var laws []domain.Laws
	var total int64

	if err := s.Repo.DB.Model(&domain.Laws{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := s.Repo.DB.Offset(offset).Limit(limit).Find(&laws).Error; err != nil {
		return nil, 0, err
	}

	return laws, total, nil
}
func (s *LawsService) GetByID(id uint) (*domain.Laws, error) {
	return s.Repo.FindByID(id)
}

func (s *LawsService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
func (s *LawsService) Update(id uint, updatedLaws *domain.Laws) (*domain.Laws, error) {
	existingLaws, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if updatedLaws.Laws != "" {
		existingLaws.Laws = updatedLaws.Laws
	}

	return existingLaws, nil
}
