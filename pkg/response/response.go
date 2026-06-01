package response

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, status int, errMsg string) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Message: "Request failed",
		Error:   errMsg,
	})
}

func ValidationError(c *fiber.Ctx, errors any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(APIResponse{
		Success: false,
		Message: "Validation failed",
		Error:   errors,
	})
}
