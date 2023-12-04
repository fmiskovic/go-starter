package user

import (
	"strconv"
	"strings"

	apiErr "github.com/fmiskovic/go-starter/internal/core/error"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/fmiskovic/go-starter/internal/core/validators"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
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

// HandleCreate creates handler func that is responsible for persisting new user entity.
// Response is UserDto json.
func (uh Handler) HandleCreate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse request body
		req := new(user.CreateRequest)
		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}

		// validate request
		if errs := uh.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// call core service
		res, err := uh.service.Create(c.Context(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrEntityCreate)).Error())
		}

		// response
		c.Status(fiber.StatusCreated)
		return toJson(c, res)
	}
}

// HandleUpdate creates handler func that is responsible for updating existing user entity.
// Response is UserDto json.
func (uh Handler) HandleUpdate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse request body
		req := new(user.UpdateRequest)
		if err := c.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}

		// validate request
		if errs := uh.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// call core service
		res, err := uh.service.Update(c.Context(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrEntityUpdate)).Error())
		}

		// response
		return toJson(c, res)
	}
}

// HandleGetById creates handler func that is responsible for getting existing user entity by its ID.
// Response is UserDto json.
func (uh Handler) HandleGetById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse query params
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithAppErr(apiErr.ErrInvalidId)).Error())
		}

		id, err := uuid.Parse(sId)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrInvalidId)).Error())
		}

		// call core service
		res, err := uh.service.GetById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrGetById)).Error())
		}

		// response
		return toJson(c, res)
	}
}

// HandleDeleteById creates handler func that is responsible for deleting existing user entity by its ID.
func (uh Handler) HandleDeleteById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse query params
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithAppErr(apiErr.ErrInvalidId)).Error())
		}

		id, err := uuid.Parse(sId)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrInvalidId)).Error())
		}

		// call core service
		err = uh.service.DeleteById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrDeleteById)).Error())
		}

		// response
		c.Status(fiber.StatusNoContent)
		return nil
	}
}

// HandleGetPage returns page of users
// HandleGetPage creates handler func that is responsible for getting page of user entities.
// Response is json representing Page of UserDtos.
func (uh Handler) HandleGetPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// parse query params
		size, err := strconv.Atoi(c.Query("size", "10"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrInvalidPageSize)).Error())
		}

		offset, err := strconv.Atoi(c.Query("offset", "0"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrInvalidPageOffset)).Error())
		}

		sort := resolveSort(c)

		pageReq := domain.Pageable{
			Size:   size,
			Offset: offset,
			Sort:   sort,
		}

		// call core service
		page, err := uh.service.GetPage(c.Context(), pageReq)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrGetPage)).Error())
		}

		// response
		return toJson(c, page)
	}
}

func toJson(c *fiber.Ctx, t interface{}) error {
	if err := c.JSON(t); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// resolveSort parses the sort parameter into a sort object.
func resolveSort(c *fiber.Ctx) domain.Sort {
	// extract sort parameters from query parameters
	sortParam := c.Query("sort", "")
	if sortParam == "" {
		return domain.NewSort()
	}
	// split the sort parameter into individual sort orderParams
	orderParams := strings.Split(sortParam, ",")

	var orders []*domain.Order
	// remove any leading or trailing spaces from each sort order
	for i := range orderParams {
		o := strings.Split(strings.TrimSpace(orderParams[i]), " ")
		order := domain.NewOrder(domain.WithProperty(o[0]), domain.WithDirection(domain.ASC))
		if len(o) == 2 {
			order.Direction = domain.Direction(o[1])
		}
		orders = append(orders, order)
	}

	return domain.NewSort(orders...)
}
