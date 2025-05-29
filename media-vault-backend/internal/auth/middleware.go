package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig defines the configuration for JWT middleware
type JWTConfig struct {
	SigningKey []byte
}

// NewAuthMiddleware creates a new JWT authentication middleware
func NewAuthMiddleware(config JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
			}
			return config.SigningKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Add user info to context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("userID", claims["sub"])
			c.Locals("roles", claims["realm_access"].(map[string]interface{})["roles"])
		}

		return c.Next()
	}
}

// RequireRole checks if the user has the required role
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles, ok := c.Locals("roles").([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No roles found in token",
			})
		}

		for _, r := range userRoles {
			if r == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}