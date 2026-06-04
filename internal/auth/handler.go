package auth

import (
	"github.com/moxhyusuf/ai-resume-analyzer/pkg/response"
	"github.com/moxhyusuf/ai-resume-analyzer/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// Register godoc
// @Summary      Register user baru
// @Description  Membuat akun baru dengan nama, email, dan password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Register payload"
// @Success      201 {object} response.APIResponse{data=AuthResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      409 {object} response.APIResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		return response.Error(c, fiber.StatusConflict, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Registration successful", result)
}

// Login godoc
// @Summary      Login user
// @Description  Login dengan email dan password, mendapatkan JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login payload"
// @Success      200 {object} response.APIResponse{data=AuthResponse}
// @Failure      401 {object} response.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Login successful", result)
}
