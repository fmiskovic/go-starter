package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fmiskovic/go-starter/internal/test"
	"github.com/fmiskovic/go-starter/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/matryer/is"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCreate(t *testing.T) {
	assert := is.New(t)

	bunDb, app := test.SetUpServer(t)

	repo := NewRepo(bunDb)
	valid := validator.New()
	app.Post("/user", HandleCreate(repo, valid))

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

				dto := &Dto{}
				err := json.NewDecoder(resBody).Decode(dto)
				assert.NoErr(err)
				assert.Equal(dto.Email, "test1@fake.com")
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

	bunDb, app := test.SetUpServer(t)

	repo := NewRepo(bunDb)
	valid := validator.New()
	app.Put("/user", HandleUpdate(repo, valid))

	tests := []struct {
		name     string
		route    string
		reqBody  *Dto
		given    func(t *testing.T) error
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given valid user request should return 200",
			route:    "/user",
			reqBody:  NewDto(Id(1), Email("test1@fake.com"), Location("Vienna")),
			wantCode: 200,
			given: func(t *testing.T) error {
				return repo.Create(context.Background(), &User{Email: "test1@fake.com"})
			},
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				dto := &Dto{}
				err := json.NewDecoder(resBody).Decode(dto)
				assert.NoErr(err)
				assert.Equal(dto.Location, "Vienna")
			},
		},
		{
			name:    "given nil user request should return 400",
			route:   "/user",
			reqBody: nil,
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
		{
			name:    "given invalid user request email should return 400",
			route:   "/user",
			reqBody: NewDto(Email("")),
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 400,
			verify:   func(t *testing.T, res *http.Response) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			err := tt.given(t)
			assert.NoErr(err)

			reqJson, err := json.Marshal(tt.reqBody)
			assert.NoErr(err)

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
	type args struct {
		repo UserRepo
	}
	tests := []struct {
		name string
		args args
		want func(c *fiber.Ctx) error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := HandleDeleteById(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("HandleDeleteById() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestHandleGetById(t *testing.T) {
	type args struct {
		repo UserRepo
	}
	tests := []struct {
		name string
		args args
		want func(c *fiber.Ctx) error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := HandleGetById(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("HandleGetById() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestHandleGetPage(t *testing.T) {
	assert := is.New(t)

	bunDb, app := test.SetUpServer(t)

	repo := NewRepo(bunDb)
	app.Get("/user", HandleGetPage(repo))

	tests := []struct {
		name     string
		route    string
		given    func(t *testing.T) error
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:  "given valid request should return 200",
			route: "/user",
			given: func(t *testing.T) error {
				return repo.Create(context.Background(), &User{Email: "test1@fake.com"})
			},
			wantCode: 200,
			verify: func(t *testing.T, res *http.Response) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			err := tt.given(t)
			assert.NoErr(err)

			req := httptest.NewRequest("GET", tt.route, nil)
			//req.Header.Add("Content-Type", "application/json")

			// when
			res, err := app.Test(req, 1000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}
