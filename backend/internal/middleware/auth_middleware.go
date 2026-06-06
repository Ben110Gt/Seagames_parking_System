package middleware

import (
	util "seagame/ticket/backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTMiddleware ตรวจสอบ JWT
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := ""

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// If no token in header, check query parameter (useful for WebSockets)
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		token, err := util.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		// ดึง user_id (string) จาก JWT แล้วเก็บใน context
		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user_id in token"})
		}
		c.Locals("user_id", userID)
		c.Locals("role", claims["role"].(string))

		return c.Next()
	}
}
func RoleMiddleware(roles ...string) fiber.Handler {
	roleMap := map[string]bool{}
	for _, r := range roles {
		roleMap[r] = true
	}

	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "role not found"})
		}

		if !roleMap[role.(string)] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
		}

		return c.Next()
	}
}

// GetUserID extracts user ID from JWT locals set by JWTMiddleware.
func GetUserID(c *fiber.Ctx) uuid.UUID {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return uuid.Nil
	}
	id, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil
	}
	return id
}
