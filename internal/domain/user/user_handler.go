package user

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
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
		if err == nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		u, err := repo.GetById(c.Context(), id)
		if err == nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
		if err == nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		err = repo.DeleteById(c.Context(), id)
		if err == nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func toJson(c *fiber.Ctx, u *User) error {
	if err := c.JSON(u); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
