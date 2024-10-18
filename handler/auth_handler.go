// handler/auth_handler.go
package handler

import (
	"rr/database"
	"rr/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Register godoc
// @Summary Register new user
// @Description Create new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func Register(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := service.RegisterUser(database.DB, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := service.LoginUser(database.DB, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	// Store token in a cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true // Cookie can't be accessed via JavaScript
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Login successful"})
}

// Logout godoc
// @Summary User logout
// @Description Invalidate user session
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /logout [post]
func Logout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour) // Expire the cookie
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Logout successful"})
}
