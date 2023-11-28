package auth

import (
	"strings"
	"time"

	apiErr "github.com/fmiskovic/go-starter/internal/adapters/api/error"
	"github.com/fmiskovic/go-starter/internal/adapters/api/validator"
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/fmiskovic/go-starter/internal/utils/password"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	repo      ports.UserRepo[uuid.UUID]
	validator validator.Validator
}

func NewHandler(repo ports.UserRepo[uuid.UUID]) AuthHandler {
	return AuthHandler{repo: repo, validator: validator.New()}
}

func (h AuthHandler) HandleSignIn() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var req = new(Request)
		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}

		if errs := h.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		user, err := h.repo.GetByUsername(c.Context(), req.Username)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrInvalidAuthReq)).Error())
		}

		if !password.CheckPasswordHash(req.Password, user.Credentials.Password) {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrInvalidAuthReq)).Error())
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"email": user.Email,
			"sub":   user.ID,
			"roles": user.Roles,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Set(fiber.HeaderAuthorization, "Bearer "+t)

		return nil
	}
}
