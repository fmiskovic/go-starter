package user

import (
	"github.com/fmiskovic/go-starter/internal/domain"
	"github.com/fmiskovic/go-starter/pkg/errorx"
	"github.com/fmiskovic/go-starter/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

// HandleCreate persists and returns new user entity
func HandleCreate(repo UserRepo, validator validator.Validator) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		r, err := parseRequestBody(c)
		if err != nil {
			return err
		}

		if errs := validator.Validate(r); len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, strings.Join(errs, " and "))
		}

		// convert request to user entity
		u := toUser(r)

		if err := repo.Create(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrUserCreate)).Error())
		}

		c.Status(fiber.StatusCreated)
		return toJson(c, u)
	}
}

func HandleUpdate(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var u = new(User)
		if err := c.BodyParser(u); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrUserUpdateReqBody)).Error())
		}

		if err := repo.Update(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrUserUpdate)).Error())
		}

		return toJson(c, u)
	}
}

func HandleGetById(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithAppErr(ErrInvalidUserId)).Error())
		}

		id, err := strconv.ParseUint(sId, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrInvalidUserId)).Error())
		}

		u, err := repo.GetById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrUserGetById)).Error())
		}

		return toJson(c, u)
	}
}

func HandleDeleteById(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithAppErr(ErrInvalidUserId)).Error())
		}

		id, err := strconv.ParseUint(sId, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrInvalidUserId)).Error())
		}

		err = repo.DeleteById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrUserDeleteById)).Error())
		}

		return nil
	}
}

func HandleGetPage(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		size, err := strconv.Atoi(c.Query("size", "10"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrInvalidPageSize)).Error())
		}

		offset, err := strconv.Atoi(c.Query("offset", "0"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrInvalidPageOffset)).Error())
		}

		sort := resolveSort(c)

		pageReq := domain.Pageable{
			Size:   size,
			Offset: offset,
			Sort:   sort,
		}
		page, err := repo.GetPage(c.Context(), pageReq)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError,
				errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrUserGetPage)).Error())
		}
		return toJson(c, page)
	}
}

func parseRequestBody(c *fiber.Ctx) (*Request, error) {
	var r = new(Request)
	if err := c.BodyParser(r); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest,
			errorx.New(errorx.WithSvcErr(err), errorx.WithAppErr(ErrUserReqBody)).Error())
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
