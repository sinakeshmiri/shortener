package httpfiber

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (ha Adapter) jwtMiddleware(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing authorization header")
	}

	token := extractTokenFromHeader(authorizationHeader)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token format")
	}

	// You can now use the 'token' variable for JWT validation or other purposes
	sub, err := ha.authc.SubjectFrom(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}
	acc,err:=ha.authc.GetAccount(sub)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}
	if c.Params("username")==acc.Username{
		return c.Next()
	}
	return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
}

func extractTokenFromHeader(authorizationHeader string) string {
	// The header value typically looks like "Bearer <token>"
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
