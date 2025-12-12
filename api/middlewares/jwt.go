package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWT protects private routes and renews access token on each request (sliding expiration).
func JWT(secret []byte, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("user_id", int64(claims["sub"].(float64)))
		if name, ok := claims["name"].(string); ok {
			c.Locals("user_name", name)
		}
		if email, ok := claims["email"].(string); ok {
			c.Locals("user_email", email)
		}
		if role, ok := claims["role"].(string); ok {
			c.Locals("user_role", role)
		}

		// Renew access token and send back via header for the client to update.
		newClaims := jwt.MapClaims{
			"sub":   claims["sub"],
			"email": claims["email"],
			"name":  claims["name"],
			"role":  claims["role"],
			"exp":   time.Now().Add(ttl).Unix(),
		}
		if newToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims).SignedString(secret); err == nil {
			c.Set("X-Access-Token", newToken)
		}

		return c.Next()
	}
}

// RequireRole checks role against allowed list.
func RequireRole(allowed ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("user_role").(string)
		for _, a := range allowed {
			if role == a {
				return c.Next()
			}
		}
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
}
