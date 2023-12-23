package auth

import (
	"log/slog"

	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/domain/security"
	"github.com/fmiskovic/go-starter/internal/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	cfg configs.AuthConfig
}

func NewMiddleware(cfg configs.AuthConfig) Middleware {
	return Middleware{cfg: cfg}
}

func (m Middleware) Authenticated() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(m.cfg.Secret)},
	})
}

func (m Middleware) AdminAuthenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if utils.IsBlank(authHeader) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		tokenString := c.Get(fiber.HeaderAuthorization)[7:]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.Secret), nil
		})

		if err != nil {
			slog.Error("parsing jwt", "error", err)
			if err == jwt.ErrSignatureInvalid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid signature",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		if !token.Valid {
			slog.Error("jwt is invalid")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		// Check if the user has the "admin" role
		roles := claims["roles"].([]interface{})
		if !containsAdminRole(roles) {
			slog.Error("ROLE_ADMIN is not present in the jwt claims")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Permission denied",
			})
		}

		c.Locals("user", token)
		return c.Next()
	}
}

func containsAdminRole(roles []interface{}) bool {
	for _, role := range roles {
		if role == security.ROLE_ADMIN {
			return true
		}
	}
	return false
}
