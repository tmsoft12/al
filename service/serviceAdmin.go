package service

import (
	"errors"
	"rr/domain"
	repository "rr/repostory"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	AdminRepo repository.AdminRepository
	Secret    string
}

func NewAdminService(repo repository.AdminRepository, secret string) *AdminService {
	return &AdminService{
		AdminRepo: repo,
		Secret:    secret,
	}
}

func (s *AdminService) Register(admin *domain.Admin) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(admin.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.Password = hashedPassword
	return s.AdminRepo.Register(admin)
}

func (s *AdminService) Login(username, password string) (string, error) {
	admin, err := s.AdminRepo.Login(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(admin.Password, []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// JWT tokenini döredýäris
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": admin.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
