package handlers

import (
	"github.com/fmiskovic/go-starter/internal/infrastructure/persistence"
	"github.com/fmiskovic/go-starter/internal/interfaces/api"
	"github.com/fmiskovic/go-starter/internal/interfaces/errorsx"
	"github.com/fmiskovic/go-starter/internal/interfaces/pagination"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HandleCreate creates handler func that is responsible for persisting new user entity.
// Response is UserDto json.
func HandleCreate(repo persistence.UserRepo, validator Validator) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req, err := parseRequestBody(c)
		if err != nil {
			return err
		}

		if errs := validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// convert request to u entity
		u := api.ToUser(req)
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()

		if err := repo.Create(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrUserCreate)).Error())
		}

		res := api.ToUserDto(u)
		c.Status(fiber.StatusCreated)
		return toJson(c, res)
	}
}

// HandleUpdate creates handler func that is responsible for updating existing user entity.
// Response is UserDto json.
func HandleUpdate(repo persistence.UserRepo, validator Validator) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req, err := parseRequestBody(c)
		if err != nil {
			return err
		}

		if errs := validator.Validate(req); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// convert request to user entity
		u := api.ToUser(req)
		u.UpdatedAt = time.Now()

		if err := repo.Update(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrUserUpdate)).Error())
		}

		return toJson(c, api.ToUserDto(u))
	}
}

// HandleGetById creates handler func that is responsible for getting existing user entity by its ID.
// Response is UserDto json.
func HandleGetById(repo persistence.UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				errorsx.New(errorsx.WithAppErr(api.ErrInvalidUserId)).Error())
		}

		id, err := strconv.ParseUint(sId, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrInvalidUserId)).Error())
		}

		u, err := repo.GetById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrUserGetById)).Error())
		}

		return toJson(c, api.ToUserDto(u))
	}
}

// HandleDeleteById creates handler func that is responsible for deleting existing user entity by its ID.
func HandleDeleteById(repo persistence.UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				errorsx.New(errorsx.WithAppErr(api.ErrInvalidUserId)).Error())
		}

		id, err := strconv.ParseUint(sId, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrInvalidUserId)).Error())
		}

		err = repo.DeleteById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrUserDeleteById)).Error())
		}

		c.Status(fiber.StatusNoContent)
		return nil
	}
}

// HandleGetPage returns page of users
// HandleGetPage creates handler func that is responsible for getting page of user entities.
// Response is json representing Page of UserDtos.
func HandleGetPage(repo persistence.UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		size, err := strconv.Atoi(c.Query("size", "10"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrInvalidPageSize)).Error())
		}

		offset, err := strconv.Atoi(c.Query("offset", "0"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrInvalidPageOffset)).Error())
		}

		sort := resolveSort(c)

		pageReq := pagination.Pageable{
			Size:   size,
			Offset: offset,
			Sort:   sort,
		}
		page, err := repo.GetPage(c.Context(), pageReq)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrUserGetPage)).Error())
		}
		return toJson(c, api.ToPageDto(page))
	}
}

func parseRequestBody(c *fiber.Ctx) (*api.UserDto, error) {
	var r = new(api.UserDto)
	if err := c.BodyParser(r); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest,
			errorsx.New(errorsx.WithSvcErr(err), errorsx.WithAppErr(api.ErrUserReqBody)).Error())
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
func resolveSort(c *fiber.Ctx) pagination.Sort {
	// extract sort parameters from query parameters
	sortParam := c.Query("sort", "")
	if sortParam == "" {
		return pagination.NewSort()
	}
	// split the sort parameter into individual sort orderParams
	orderParams := strings.Split(sortParam, ",")

	var orders []*pagination.Order
	// remove any leading or trailing spaces from each sort order
	for i := range orderParams {
		o := strings.Split(strings.TrimSpace(orderParams[i]), " ")
		order := pagination.NewOrder(pagination.WithProperty(o[0]), pagination.WithDirection(pagination.ASC))
		if len(o) == 2 {
			order.Direction = pagination.Direction(o[1])
		}
		orders = append(orders, order)
	}

	return pagination.NewSort(orders...)
}
