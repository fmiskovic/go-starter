package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/core/services"
	"github.com/fmiskovic/go-starter/internal/utils/testx"
	"github.com/google/uuid"

	"github.com/matryer/is"
)

func TestHandleCreate(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Post("/user", handler.HandleCreate())

	tests := []struct {
		name     string
		route    string
		reqBody  []byte
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given valid create request should return 201",
			route:    "/user",
			reqBody:  []byte("{\"username\":\"test1\",\"password\":\"Password1234!\",\"email\":\"test1@fake.com\"}"),
			wantCode: 201,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				createRes := &user.CreateResponse{}
				err := json.NewDecoder(resBody).Decode(createRes)
				assert.NoErr(err)
				assert.Equal(createRes.Email, "test1@fake.com")
			},
		},
		{
			name:     "given empty create request should return 400",
			route:    "/user",
			reqBody:  []byte(""),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given invalid email should return 400",
			route:    "/user",
			reqBody:  []byte("{\"username\":\"Password1234!\",\"password\":\"test1\",\"email\":\"\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tt.route, bytes.NewReader(tt.reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 20000)
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleUpdate(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Put("/user", handler.HandleUpdate())

	tests := []struct {
		name     string
		route    string
		reqBody  []byte
		verify   func(t *testing.T, res *http.Response)
		wantCode int
	}{
		{
			name:     "given valid update request should return 200",
			route:    "/user",
			reqBody:  []byte("{\"id\":\"220cea28-b2b0-4051-9eb6-9a99e451af01\",\"email\":\"test1@fake.com\", \"location\":\"Vienna\"}"),
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				updateRes := &user.UpdateResponse{}
				err := json.NewDecoder(resBody).Decode(updateRes)
				assert.NoErr(err)
				assert.Equal(updateRes.Location, "Vienna")
			},
		},
		{
			name:     "given empty update request should return 400",
			route:    "/user",
			reqBody:  []byte("{}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given invalid email should return 400",
			route:    "/user",
			reqBody:  []byte("{\"id\":\"220cea28-b2b0-4051-9eb6-9a99e451af01\",\"email\":\"\"}"),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given invalid id should return 200",
			route:    "/user",
			reqBody:  []byte("{\"id\":\"333cea28-b2b0-4051-9eb6-9a99e451af01\",\"email\":\"test1@fake.com\"}"),
			wantCode: 200,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			req := httptest.NewRequest("PUT", tt.route, bytes.NewReader(tt.reqBody))
			req.Header.Add("Content-Type", "application/json")

			// when
			res, err := app.Test(req, 5000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleDeleteById(t *testing.T) {
	assert := is.NewRelaxed(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Delete("/user/:id", handler.HandleDeleteById())

	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		verify   func(id string, t *testing.T)
	}{
		{
			name:     "given user id should return 204 and delete user",
			wantCode: 204,
			args:     args{id: "220cea28-b2b0-4051-9eb6-9a99e451af01"},
			verify: func(id string, t *testing.T) {
				u, err := repo.GetById(context.Background(), uuid.MustParse(id))
				assert.True(err != nil)
				assert.True(u == nil)
			},
		},
		{
			name:     "given non-existing user id should return 204",
			args:     args{id: "333cea28-b2b0-4051-9eb6-9a99e451af01"},
			wantCode: 204,
			verify: func(id string, t *testing.T) {
			},
		},
		{
			name:     "given empty user id should return 404",
			args:     args{id: ""},
			wantCode: 404,
			verify: func(id string, t *testing.T) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			req := httptest.NewRequest("DELETE", fmt.Sprintf("%s/%s", "/user", tt.args.id), nil)
			// when
			res, err := app.Test(req, 5000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(tt.args.id, t)
		})
	}
}

func TestHandleGetById(t *testing.T) {
	assert := is.NewRelaxed(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Get("/user/:id", handler.HandleGetById())

	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given user id should return 200 and user dto",
			wantCode: 200,
			args:     args{id: "220cea28-b2b0-4051-9eb6-9a99e451af01"},
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				userDto := &user.Dto{}
				err := json.NewDecoder(resBody).Decode(userDto)
				assert.NoErr(err)
				assert.Equal(userDto.Email, "john@smith.com")
				assert.Equal(userDto.Gender.Numberfy(), user.MALE)
			},
		},
		{
			name:     "given non-existing user id should return 404",
			args:     args{id: "333cea28-b2b0-4051-9eb6-9a99e451af01"},
			wantCode: 404,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
		{
			name:     "given empty user id should return 404",
			args:     args{id: ""},
			wantCode: 404,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", "user", tt.args.id), nil)

			// when
			res, err := app.Test(req, 5000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleGetPage(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Get("/user", handler.HandleGetPage())

	tests := []struct {
		name     string
		route    string
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given empty pageable should return 200",
			route:    "/user",
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				var pageDto domain.Page[user.Dto]
				err := json.NewDecoder(resBody).Decode(&pageDto)
				assert.NoErr(err)
				assert.True(pageDto.TotalElements > 0)
				assert.Equal(pageDto.Elements[0].Email, "john@smith.com")
			},
		},
		{
			name:     "given pageable should return 200",
			route:    "/user?size=10&offset=0&sort=email%20ASC",
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				var pageDto domain.Page[user.Dto]
				err := json.NewDecoder(resBody).Decode(&pageDto)
				assert.NoErr(err)
				assert.True(pageDto.TotalElements > 0)
				assert.Equal(pageDto.Elements[1].Email, "john@doe.com")
			},
		},
		{
			name:     "given pageable with offset 5 should return 200 and no elements",
			route:    "/user?offset=5&sort=email%20ASC",
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				var pageDto domain.Page[user.Dto]
				err := json.NewDecoder(resBody).Decode(&pageDto)
				assert.NoErr(err)
				assert.True(len(pageDto.Elements) == 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			req := httptest.NewRequest("GET", tt.route, nil)
			// when
			res, err := app.Test(req, 5000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleUserRoles(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Post("/user/roles", handler.HandleUserRoles())

	tests := []struct {
		name     string
		reqBody  []byte
		wantCode int
	}{
		{
			name:     "given valid add request should return 201",
			reqBody:  []byte("{\"id\":\"220cea28-b2b0-4051-9eb6-9a99e451af02\",\"command\":\"ADD\",\"roles\":[\"ROLE_ADMIN\"]}"),
			wantCode: 201,
		},
		{
			name:     "given valid delete request should return 204",
			reqBody:  []byte("{\"id\":\"220cea28-b2b0-4051-9eb6-9a99e451af02\",\"command\":\"DELETE\",\"roles\":[\"ROLE_ADMIN\"]}"),
			wantCode: 204,
		},
		{
			name:     "given empty request should return 400",
			reqBody:  []byte(""),
			wantCode: 400,
		},
		{
			name:     "given invalid command should return 400",
			reqBody:  []byte("{\"id\":\"220cea28-b2b0-4051-9eb6-9a99e451af02\",\"command\":\"INVALID\",\"roles\":[\"ROLE_ADMIN\"]}"),
			wantCode: 400,
		},
		{
			name:     "given invalid id should return 422",
			reqBody:  []byte("{\"id\":\"invalid\",\"command\":\"ADD\",\"roles\":[\"ROLE_ADMIN\"]}"),
			wantCode: 422,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/user/roles", bytes.NewReader(tt.reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 20000)
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
		})
	}
}

func TestHandleEnableDisable(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	service := services.NewUserService(repo, configs.NewAuthConfig())
	handler := NewHandler(service)
	app.Post("/user/:id/enabledisable", handler.HandleEnableDisable())

	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		verify   func(t *testing.T, id string)
		wantCode int
	}{
		{
			name: "given valid id should enable user",
			args: args{id: "220cea28-b2b0-4051-9eb6-9a99e451af02"},
			verify: func(t *testing.T, id string) {
				u, err := repo.GetById(context.Background(), uuid.MustParse(id))
				if err != nil {
					t.Errorf("failed to get user by id, error: %s", err.Error())
				}
				assert.True(u.Enabled)
			},
			wantCode: 200,
		},
		{
			name: "given valid id should disable user",
			args: args{id: "220cea28-b2b0-4051-9eb6-9a99e451af02"},
			verify: func(t *testing.T, id string) {
				u, err := repo.GetById(context.Background(), uuid.MustParse(id))
				if err != nil {
					t.Errorf("failed to get user by id, error: %s", err.Error())
				}
				assert.True(!u.Enabled)
			},
			wantCode: 200,
		},
		{
			name:     "given invalid id should return 422",
			args:     args{id: "333cea28-b2b0-4051-9eb6-9a99e451af02"},
			verify:   func(t *testing.T, id string) {},
			wantCode: 422,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", fmt.Sprintf("/user/%s/enabledisable", tt.args.id), nil)
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 20000)
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, tt.args.id)
		})
	}
}
