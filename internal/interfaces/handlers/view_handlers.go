package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
	"net/http"
)

func HandleHome(c *fiber.Ctx) error {
	return c.Render("home/index", fiber.Map{})
}

func HandleAbout(c *fiber.Ctx) error {
	return c.Render("home/about", fiber.Map{})
}

func HandleFlash(c *fiber.Ctx) error {
	context := fiber.Map{
		"systemMessage": "a flash message for you user",
	}
	return flash.WithData(c, context).RedirectBack("/")
}

func NotFoundMiddleware(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).Render("error/404", nil)
}

func FlashMiddleware(c *fiber.Ctx) error {
	c.Locals("flash", flash.Get(c))
	return c.Next()
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	return c.Render("error/500", nil)
}
