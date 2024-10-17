package handler

import (
	"rr/domain"
	"rr/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	AdminService *service.AdminService
}

func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{
		AdminService: service,
	}
}

func (h *AdminHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password []byte `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	admin := &domain.Admin{
		Username: req.Username,
		Password: req.Password,
	}

	if err := h.AdminService.Register(admin); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "registered successfully"})
}

func (h *AdminHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.AdminService.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Tokeni cookie-da saklaýarys
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"message": "logged in successfully"})
}
