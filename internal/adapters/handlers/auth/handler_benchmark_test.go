package auth

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/services"
	"github.com/fmiskovic/go-starter/internal/utils/testx"
	"github.com/matryer/is"
)

func BenchmarkHandleSignIn(b *testing.B) {
	if testing.Short() {
		return
	}
	assert := is.New(b)

	ts, err := testx.SetUpServer()
	if err != nil {
		b.Errorf("failed to run test server: %v", err)
	}
	defer ts.TestDb.Shutdown()

	repo := repos.NewUserRepo(ts.TestDb.BunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	ts.App.Post("/auth/login", handler.HandleSignIn())

	for n := 0; n < b.N; n++ {
		body := []byte("{\"username\":\"username1\",\"password\":\"password1\"}")
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")

		res, err := ts.App.Test(req, 20000)
		assert.NoErr(err)
		assert.Equal(res.StatusCode, 200)
	}
}

func BenchmarkHandleSignUp(b *testing.B) {
	if testing.Short() {
		return
	}
	assert := is.New(b)

	ts, err := testx.SetUpServer()
	if err != nil {
		b.Errorf("failed to run test server: %v", err)
	}
	defer ts.TestDb.Shutdown()

	repo := repos.NewUserRepo(ts.TestDb.BunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	ts.App.Post("/auth/register", handler.HandleSignUp())

	for n := 0; n < b.N; n++ {
		body := []byte("{\"username\":\"test1\",\"password\":\"Password1234!\",\"email\":\"test1@fake.com\"}")
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")

		res, err := ts.App.Test(req, 20000)
		assert.NoErr(err)
		assert.Equal(res.StatusCode, 201)
	}
}
