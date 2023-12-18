package user

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

func BenchmarkHandleUpdate(b *testing.B) {
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
	ts.App.Put("/user", handler.HandleUpdate())

	for n := 0; n < b.N; n++ {
		body := []byte("{\"id\":\"220cea28-b2b0-4051-9eb6-9a99e451af01\",\"email\":\"test1@fake.com\", \"location\":\"Vienna\"}")
		req := httptest.NewRequest("PUT", "/user", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")

		res, err := ts.App.Test(req, 20000)

		assert.NoErr(err)
		assert.Equal(res.StatusCode, 200)
	}
}
