package auth

import (
	"strings"

	apiErr "github.com/fmiskovic/go-starter/internal/core/error"

	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/fmiskovic/go-starter/internal/core/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service   ports.UserService[uuid.UUID]
	validator validators.Validator
}

func NewHandler(service ports.UserService[uuid.UUID]) Handler {
	return Handler{
		service:   service,
		validator: validators.New(),
	}
}

// HandleSingIn is used to authenticate user.
func (h Handler) HandleSignIn() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse request body
		var req = new(user.SignInRequest)
		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}

		// validate request
		if errs := h.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// call core service
		res, err := h.service.SingIn(c.Context(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrInvalidAuthReq)).Error())
		}

		// response
		c.Set(fiber.HeaderAuthorization, "Bearer "+res.Token)
		return c.JSON(res)
	}
}

// HandleSignUp is used to register new user.
func (h Handler) HandleSignUp() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse request body
		var req = new(user.CreateRequest)
		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}

		// validate request
		if errs := h.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// call core service
		res, err := h.service.SingUp(c.Context(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrSignUp)).Error())
		}

		// response
		c.Status(fiber.StatusCreated)
		return c.JSON(res)
	}
}

// HandleChangePassword change user password.
func (h Handler) HandleChangePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse request body
		var req = new(user.ChangePasswordRequest)
		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}

		// validate request
		if errs := h.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// call core service
		if err := h.service.ChangePassword(c.Context(), req); err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, apiErr.New(apiErr.WithSvcErr(err)).Error())
		}

		// response
		c.Status(fiber.StatusNoContent)
		return nil
	}
}

// HandleConfirmEmail confirms user email address.
func (h Handler) HandleConfirmEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse query params
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithAppErr(apiErr.ErrInvalidId)).Error())
		}
		code := c.Params("code", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithAppErr(apiErr.ErrInvalidCode)).Error())
		}

		req := new(user.ConfirmEmailRequest)
		req.ID = sId
		req.Code = code

		// call core service
		err := h.service.ConfirmEmail(c.Context(), *req)
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, apiErr.New(apiErr.WithSvcErr(err)).Error())
		}

		// response
		c.Status(fiber.StatusNoContent)
		return nil
	}
}

// HandleSignOut logout user.
func (h Handler) HandleSignOut() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user", nil)
		c.Set(fiber.HeaderAuthorization, "Bearer ")
		return nil
	}
}
