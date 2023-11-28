package user

import (
	"strconv"
	"strings"
	"time"

	apiErr "github.com/fmiskovic/go-starter/internal/adapters/api/error"
	"github.com/fmiskovic/go-starter/internal/adapters/api/validator"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo      ports.UserRepo[uuid.UUID]
	validator validator.Validator
}

func NewHandler(repo ports.UserRepo[uuid.UUID]) Handler {
	return Handler{
		repo:      repo,
		validator: validator.New(),
	}
}

// HandleCreate creates handler func that is responsible for persisting new user entity.
// Response is UserDto json.
func (uh Handler) HandleCreate() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req, err := parseRequestBody(c)
		if err != nil {
			return err
		}

		if errs := uh.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// convert request to u entity
		u, err := ToUser(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()

		if err := uh.repo.Create(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrEntityCreate)).Error())
		}

		res := ToDto(u)
		c.Status(fiber.StatusCreated)
		return toJson(c, res)
	}
}

// HandleUpdate creates handler func that is responsible for updating existing user entity.
// Response is UserDto json.
func (uh Handler) HandleUpdate() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req, err := parseRequestBody(c)
		if err != nil {
			return err
		}

		if errs := uh.validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// convert request to user entity
		u, err := ToUser(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
		}
		u.UpdatedAt = time.Now()

		if err := uh.repo.Update(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrEntityUpdate)).Error())
		}

		return toJson(c, ToDto(u))
	}
}

// HandleGetById creates handler func that is responsible for getting existing user entity by its ID.
// Response is UserDto json.
func (uh Handler) HandleGetById() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		u, err := uh.repo.GetById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrGetById)).Error())
		}

		return toJson(c, ToDto(u))
	}
}

// HandleDeleteById creates handler func that is responsible for deleting existing user entity by its ID.
func (uh Handler) HandleDeleteById() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		err = uh.repo.DeleteById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrDeleteById)).Error())
		}

		c.Status(fiber.StatusNoContent)
		return nil
	}
}

// HandleGetPage returns page of users
// HandleGetPage creates handler func that is responsible for getting page of user entities.
// Response is json representing Page of UserDtos.
func (uh Handler) HandleGetPage() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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
		page, err := uh.repo.GetPage(c.Context(), pageReq)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrGetPage)).Error())
		}
		return toJson(c, ToPageDto(page))
	}
}

func parseRequestBody(c *fiber.Ctx) (*Dto, error) {
	var r = new(Dto)
	if err := c.BodyParser(r); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest,
			apiErr.New(apiErr.WithSvcErr(err), apiErr.WithAppErr(apiErr.ErrParseReqBody)).Error())
	}
	return r, nil
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
