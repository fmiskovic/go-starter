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
	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/utils/testx"
	"github.com/google/uuid"

	"github.com/matryer/is"
)

func TestHandleCreate(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	handler := NewHandler(repo)
	app.Post("/user", handler.HandleCreate())

	tests := []struct {
		name     string
		route    string
		reqBody  *Dto
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given valid user request should return 201",
			route:    "/user",
			reqBody:  NewDto(Email("test1@fake.com")),
			wantCode: 201,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				userDto := &Dto{}
				err := json.NewDecoder(resBody).Decode(userDto)
				assert.NoErr(err)
				assert.Equal(userDto.Email, "test1@fake.com")
			},
		},
		{
			name:     "given nil user request should return 400",
			route:    "/user",
			reqBody:  nil,
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:     "given invalid user request email should return 400",
			route:    "/user",
			reqBody:  NewDto(Email("")),
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqJson, err := json.Marshal(tt.reqBody)
			assert.NoErr(err)

			req := httptest.NewRequest("POST", tt.route, bytes.NewReader(reqJson))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 1000)
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
	handler := NewHandler(repo)
	app.Put("/user", handler.HandleUpdate())

	tests := []struct {
		name     string
		route    string
		reqBody  *Dto
		given    func(t *testing.T) (uuid.UUID, error)
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given valid user request should return 200",
			route:    "/user",
			reqBody:  NewDto(Email("test1@fake.com"), Location("Vienna")),
			wantCode: 200,
			given: func(t *testing.T) (uuid.UUID, error) {
				u := user.New(user.Email("test1@fake.com"))
				err := repo.Create(context.Background(), u)
				return u.ID, err
			},
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				userDto := &Dto{}
				err := json.NewDecoder(resBody).Decode(userDto)
				assert.NoErr(err)
				assert.Equal(userDto.Location, "Vienna")
			},
		},
		{
			name:    "given nil user request should return 400",
			route:   "/user",
			reqBody: nil,
			given: func(t *testing.T) (uuid.UUID, error) {
				return uuid.New(), nil
			},
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:    "given invalid user request email should return 400",
			route:   "/user",
			reqBody: NewDto(Email("")),
			given: func(t *testing.T) (uuid.UUID, error) {
				return uuid.New(), nil
			},
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			id, err := tt.given(t)
			assert.NoErr(err)

			var reqJson []byte
			userDto := tt.reqBody
			if userDto != nil {
				userDto.ID = id.String()
				reqJson, err = json.Marshal(userDto)
				assert.NoErr(err)
			}

			req := httptest.NewRequest("PUT", tt.route, bytes.NewReader(reqJson))
			req.Header.Add("Content-Type", "application/json")

			// when
			res, err := app.Test(req, 1000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleDeleteById(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	handler := NewHandler(repo)
	app.Delete("/user/:id", handler.HandleDeleteById())

	tests := []struct {
		name     string
		given    func(t *testing.T) (string, error)
		wantCode int
		verify   func(id string, t *testing.T)
	}{
		{
			name: "given user id should return 204 and delete user",
			given: func(t *testing.T) (string, error) {
				u := user.New(user.Email("test1@fake.com"))
				err := repo.Create(context.Background(), u)
				return u.ID.String(), err
			},
			wantCode: 204,
			verify: func(id string, t *testing.T) {
				u, err := repo.GetById(context.Background(), uuid.MustParse(id))
				assert.True(err != nil)
				assert.True(u == nil)
			},
		},
		{
			name: "given non-existing user id should return 204",
			given: func(t *testing.T) (string, error) {
				id := uuid.New()
				return id.String(), nil
			},
			wantCode: 204,
			verify: func(id string, t *testing.T) {
			},
		},
		{
			name: "given zero user id should return 400",
			given: func(t *testing.T) (string, error) {
				return "0", nil
			},
			wantCode: 400,
			verify: func(id string, t *testing.T) {
			},
		},
		{
			name: "given invalid user id should return 400",
			given: func(t *testing.T) (string, error) {
				return "s", nil
			},
			wantCode: 400,
			verify: func(id string, t *testing.T) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			sId, err := tt.given(t)
			assert.NoErr(err)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("%s/%s", "/user", sId), nil)

			// when
			res, err := app.Test(req, 1000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(sId, t)
		})
	}
}

func TestHandleGetById(t *testing.T) {
	assert := is.New(t)

	bunDb, app := testx.SetUpServer(t)

	repo := repos.NewUserRepo(bunDb)
	handler := NewHandler(repo)
	app.Get("/user/:id", handler.HandleGetById())

	tests := []struct {
		name     string
		given    func(t *testing.T) (string, error)
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name: "given user id should return 200 and user dto",
			given: func(t *testing.T) (string, error) {
				u := user.New(user.Email("test1@fake.com"), user.Sex(user.MALE))
				err := repo.Create(context.Background(), u)
				return u.ID.String(), err
			},
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				userDto := &Dto{}
				err := json.NewDecoder(resBody).Decode(userDto)
				assert.NoErr(err)
				assert.Equal(userDto.Email, "test1@fake.com")
				assert.Equal(userDto.Gender.Numberfy(), user.MALE)
			},
		},
		{
			name: "given non-existing user id should return 404",
			given: func(t *testing.T) (string, error) {
				return uuid.NewString(), nil
			},
			wantCode: 404,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
		{
			name: "given zero user id should return 400",
			given: func(t *testing.T) (string, error) {
				return "0", nil
			},
			wantCode: 400,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
		{
			name: "given invalid user id should return 400",
			given: func(t *testing.T) (string, error) {
				return "s", nil
			},
			wantCode: 400,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			sId, err := tt.given(t)
			assert.NoErr(err)

			req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", "user", sId), nil)

			// when
			res, err := app.Test(req, 1000)
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
	handler := NewHandler(repo)
	app.Get("/user", handler.HandleGetPage())

	tests := []struct {
		name     string
		route    string
		given    func(t *testing.T) error
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:  "given empty pageable should return 200",
			route: "/user",
			given: func(t *testing.T) error {
				u := user.New(user.Email("test1@fake.com"))
				return repo.Create(context.Background(), u)
			},
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				var pageDto domain.Page[Dto]
				err := json.NewDecoder(resBody).Decode(&pageDto)
				assert.NoErr(err)
				assert.True(pageDto.TotalElements > 0)
				assert.Equal(pageDto.Elements[0].Email, "test1@fake.com")
			},
		},
		{
			name:  "given pageable should return 200",
			route: "/user?size=10&offset=0&sort=email%20ASC",
			given: func(t *testing.T) error {
				u := user.New(user.Email("test2@fake.com"))
				return repo.Create(context.Background(), u)
			},
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				var pageDto domain.Page[Dto]
				err := json.NewDecoder(resBody).Decode(&pageDto)
				assert.NoErr(err)
				assert.True(pageDto.TotalElements > 0)
				assert.Equal(pageDto.Elements[1].Email, "test2@fake.com")
			},
		},
		{
			name:  "given pageable with offset 5 should return 200 and no elements",
			route: "/user?offset=5&sort=email%20ASC",
			given: func(t *testing.T) error {
				u := user.New(user.Email("test3@fake.com"))
				return repo.Create(context.Background(), u)
			},
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				var pageDto domain.Page[Dto]
				err := json.NewDecoder(resBody).Decode(&pageDto)
				assert.NoErr(err)
				assert.True(len(pageDto.Elements) == 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			err := tt.given(t)
			assert.NoErr(err)

			req := httptest.NewRequest("GET", tt.route, nil)

			// when
			res, err := app.Test(req, 1000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}
