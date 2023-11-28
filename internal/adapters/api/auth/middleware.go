package auth

import (
	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/contrib/jwt"
)

// Middleware is auth middleware that should be used in protected endpoints.
func Middleware(app *fiber.App, secret string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		app.Use(jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		}))
		return nil
	}
}
