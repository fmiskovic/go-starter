package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/core/services"
	"github.com/fmiskovic/go-starter/internal/utils/testx"
	"github.com/matryer/is"
)

func TestHandler_HandleSignIn(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Post("/auth/login", handler.HandleSignIn())

	tests := []struct {
		name     string
		reqBody  []byte
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given valid credentials should return 200 and token",
			reqBody:  []byte("{\"username\":\"username1\",\"password\":\"password1\"}"),
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				signInRes := &user.SignInResponse{}
				err := json.NewDecoder(resBody).Decode(signInRes)
				assert.NoErr(err)
				assert.True(signInRes.Token != "")
			},
		},
		{
			name:     "given invalid password should return 400",
			reqBody:  []byte("{\"username\":\"username1\",\"password\":\"invalid\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given non-exisitng username should return 400",
			reqBody:  []byte("{\"username\":\"non-exising\",\"password\":\"password1\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(tt.reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 1000)
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}
