package user

import (
	"errors"
	"github.com/fmiskovic/go-starter/internal/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func HandleCreate(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var u = new(User)
		if err := c.BodyParser(u); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := repo.Create(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return toJson(c, u)
	}
}

func HandleUpdate(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var u = new(User)
		if err := c.BodyParser(u); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := repo.Update(c.Context(), u); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return toJson(c, u)
	}
}

func HandleGetById(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
		}

		id, err := strconv.ParseUint(sId, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		u, err := repo.GetById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return toJson(c, u)
	}
}

func HandleDeleteById(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
		}

		id, err := strconv.ParseUint(sId, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		err = repo.DeleteById(c.Context(), id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func HandleGetPage(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		size, err := strconv.Atoi(c.Query("size", "10"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errors.Join(errors.New("invalid size number"), err).Error())
		}

		offset, err := strconv.Atoi(c.Query("offset", "0"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				errors.Join(errors.New("invalid offset number"), err).Error())
		}

		sort := resolveSort(c)

		pageReq := domain.Pageable{
			Size:   size,
			Offset: offset,
			Sort:   sort,
		}
		page, err := repo.GetPage(c.Context(), pageReq)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
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
