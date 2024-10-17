package service

import (
	"errors"
	"rr/domain"
	"rr/repostrory"
)

type EmployerService struct {
	Repo *repostrory.EmployerRepo
}

func (s *EmployerService) Create(t *domain.Employer) error {

	if t.Name == "" {
		t.Name = "Ady"
	}
	if t.Major == "" {
		t.Major = "Wezipesi"
	}
	if t.Surname == "" {
		t.Surname = "Familyasy"
	}
	if t.Image == "" {
		return errors.New("Surat hokman bolmaly")
	}
	return s.Repo.Create(t)
}
func (s *EmployerService) GetAll() ([]domain.Employer, error) {
	return s.Repo.FindAll()
}
func (s *EmployerService) GetPaginated(page int, limit int) ([]domain.Employer, int64, error) {
	var employer []domain.Employer
	var total int64

	if err := s.Repo.DB.Model(&domain.Employer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := s.Repo.DB.Offset(offset).Limit(limit).Find(&employer).Error; err != nil {
		return nil, 0, err
	}

	return employer, total, nil
}
func (s *EmployerService) GetByID(id uint) (*domain.Employer, error) {
	return s.Repo.FindByID(id)
}

func (s *EmployerService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
func (s *EmployerService) Update(id uint, updatedBanner *domain.Employer) (*domain.Employer, error) {
	existingEmployer, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if updatedBanner.Name != "" {
		existingEmployer.Name = updatedBanner.Name
	}
	if updatedBanner.Image != "" {
		existingEmployer.Image = updatedBanner.Image
	}
	if updatedBanner.Major != "" {
		existingEmployer.Major = updatedBanner.Major
	}
	if updatedBanner.Surname != "" {
		existingEmployer.Surname = updatedBanner.Surname
	}

	if err := s.Repo.Update(id, existingEmployer); err != nil {
		return nil, err
	}

	return existingEmployer, nil
}
