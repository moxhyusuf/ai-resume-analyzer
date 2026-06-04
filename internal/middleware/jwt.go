package middleware

import (
	"strings"

	"github.com/moxhyusuf/ai-resume-analyzer/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTProtected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Error(c, fiber.StatusUnauthorized, "Authorization header required")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token claims")
		}

		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid user ID in token")
		}

		c.Locals("user_id", userID)
		c.Locals("email", claims["email"])

		return c.Next()
	}
}
