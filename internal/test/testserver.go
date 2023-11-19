package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/uptrace/bun"
	"testing"
)

func SetUpServer(t *testing.T) (*bun.DB, *fiber.App) {
	t.Helper()

	_, _, bunDb := SetUpDb(t)

	app := fiber.New()
	app.Use(recover.New())

	return bunDb, app
}
