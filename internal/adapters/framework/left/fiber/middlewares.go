package httpfiber

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("authn")

func (ha Adapter) authnMiddleware(c *fiber.Ctx) error {
	_, span := tracer.Start(c.Context(), "authn", oteltrace.WithAttributes())
	defer span.End()
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing authorization header")
	}

	token := extractTokenFromHeader(authorizationHeader)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token format,no bearer")
	}

	// You can now use the 'token' variable for JWT validation or other purposes
	sub, err := ha.authc.SubjectFrom(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token " + err.Error())
	}
	//sub
	c.Locals("username", sub)
	return c.Next()

}

func extractTokenFromHeader(authorizationHeader string) string {
	// The header value typically looks like "Bearer <token>"
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
