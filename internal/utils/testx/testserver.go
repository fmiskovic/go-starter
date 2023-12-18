package testx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type TestServer struct {
	TestDb TestDB
	App    *fiber.App
}

// SetUpServer helps to set up test Server.
func SetUpServer() (*TestServer, error) {
	testDb, err := SetUpDb()
	if err != nil {
		return nil, err
	}

	app := fiber.New()
	app.Use(recover.New())

	return &TestServer{
		TestDb: *testDb,
		App:    app,
	}, nil
}
