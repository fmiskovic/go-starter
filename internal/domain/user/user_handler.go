package user

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func HandleCreate(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var u = new(User)
		if err := c.BodyParser(u); err != nil {
			return err
		}

		return repo.Create(c.Context(), u)
	}
}

func HandleUpdate(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var u = new(User)
		if err := c.BodyParser(u); err != nil {
			return err
		}

		return repo.Update(c.Context(), u)
	}
}

func HandleGetById(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			//TODO error
		}

		id, err := strconv.ParseInt(sId, 10, 64)
		if err == nil {
			//TODO error
		}

		user, err := repo.GetById(c.Context(), id)
		if err == nil {
			//TODO error
		}

		return c.JSON(user)
	}
}

func HandleDeleteById(repo UserRepo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sId := c.Params("id", "0")
		if sId == "0" {
			//TODO error
		}

		id, err := strconv.ParseInt(sId, 10, 64)
		if err == nil {
			//TODO error
		}
		return repo.DeleteById(c.Context(), id)
	}
}
