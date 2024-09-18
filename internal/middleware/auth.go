package middleware

import (
	"amg/internal/identity/accesscontrol"
	"amg/internal/identity/user"
	"amg/pkg/util/jwt"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware to check if the user has a valid JWT
func JWTProtected(secret string, service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		tokenStr := authHeader[len("Bearer "):]

		claims, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		}

		// Check if the token is blacklisted
		isBlacklisted, err := service.IsTokenBlacklisted(context.Background(), tokenStr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error while checking token",
			})
		}

		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is blacklisted",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// Middleware to check if the user has the required role
func RequireRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string) // Get the user's role from context

		// Check if the user's role is in the list of allowed roles
		roleAllowed := false
		for _, allowedRole := range requiredRoles {
			if role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Insufficient permissions.",
			})
		}

		return c.Next()
	}
}

// RequirePermission checks if the user has the required permission
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		if !accesscontrol.HasPermission(role, permission) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to access this resource",
			})
		}

		// if !accesscontrol.HasTaskPermission(role, permission) {
		// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		// 		"error": "You do not have permission to access this resource",
		// 	})
		// }

		return c.Next()
	}
}
