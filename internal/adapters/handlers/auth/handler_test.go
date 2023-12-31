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

func TestHandleSignIn(t *testing.T) {
	if testing.Short() {
		return
	}
	assert := is.New(t)

	ts, err := testx.SetUpServer()
	if err != nil {
		t.Errorf("failed to run test server: %v", err)
	}
	defer ts.TestDb.Shutdown()

	repo := repos.NewUserRepo(ts.TestDb.BunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	ts.App.Post("/auth/login", handler.HandleSignIn())

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
			name:     "given non-existing username should return 400",
			reqBody:  []byte("{\"username\":\"non-existing\",\"password\":\"password1\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(tt.reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := ts.App.Test(req, 20000)
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleSignUp(t *testing.T) {
	if testing.Short() {
		return
	}
	assert := is.New(t)

	ts, err := testx.SetUpServer()
	if err != nil {
		t.Errorf("failed to run test server: %v", err)
	}
	defer ts.TestDb.Shutdown()

	repo := repos.NewUserRepo(ts.TestDb.BunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	ts.App.Post("/auth/register", handler.HandleSignUp())

	tests := []struct {
		name     string
		reqBody  []byte
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given valid singup request should return 201",
			reqBody:  []byte("{\"username\":\"test1\",\"password\":\"Password1234!\",\"email\":\"test1@fake.com\"}"),
			wantCode: 201,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				signUpRes := &user.SignUpResponse{}
				err := json.NewDecoder(resBody).Decode(signUpRes)
				assert.NoErr(err)
				assert.True(signUpRes.ID != "")
			},
		},
		{
			name:     "given empty singup request should return 400",
			reqBody:  []byte(""),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given invalid email should return 400",
			reqBody:  []byte("{\"username\":\"username11\",\"password\":\"password11\",\"email\":\"\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given invalid username should return 400",
			reqBody:  []byte("{\"username\":\"\",\"password\":\"password11\",\"email\":\"fake@test.com\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given invalid password should return 400",
			reqBody:  []byte("{\"username\":\"username12\",\"password\":\"123\",\"email\":\"fake@test.com\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given existing username should return 422",
			reqBody:  []byte("{\"username\":\"test1\",\"password\":\"Password1234\",\"email\":\"test1@test.com\"}"),
			wantCode: 422,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(tt.reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := ts.App.Test(req, 20000)
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleChangePassword(t *testing.T) {
	if testing.Short() {
		return
	}
	assert := is.New(t)

	ts, err := testx.SetUpServer()
	if err != nil {
		t.Errorf("failed to run test server: %v", err)
	}
	defer ts.TestDb.Shutdown()

	repo := repos.NewUserRepo(ts.TestDb.BunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	ts.App.Post("/auth/password", handler.HandleChangePassword())

	tests := []struct {
		name     string
		reqBody  []byte
		wantCode int
	}{
		// {
		// 	name:     "given valid request should change user password",
		// 	reqBody:  []byte("{\"username\":\"username1\",\"oldPassword\":\"password1\",\"newPassword\":\"password111\"}"),
		// 	wantCode: 204,
		// },
		{
			name:     "given invalid oldPassword should return 422",
			reqBody:  []byte("{\"username\":\"username1\",\"oldPassword\":\"password231\",\"newPassword\":\"password111\"}"),
			wantCode: 422,
		},
		{
			name:     "given invalid newPassword should return 400",
			reqBody:  []byte("{\"username\":\"username1\",\"oldPassword\":\"password1\",\"newPassword\":\"pass\"}"),
			wantCode: 400,
		},
		{
			name:     "given invalid username should return 422",
			reqBody:  []byte("{\"username\":\"username1111\",\"oldPassword\":\"password1\",\"newPassword\":\"password111\"}"),
			wantCode: 422,
		},
		{
			name:     "given empty reques should return 400",
			reqBody:  []byte(""),
			wantCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/password", bytes.NewReader(tt.reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := ts.App.Test(req, 30000)
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
		})
	}
}
