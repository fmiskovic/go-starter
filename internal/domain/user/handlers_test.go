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

	"github.com/fmiskovic/go-starter/internal/domain"
	"github.com/fmiskovic/go-starter/internal/test"
	"github.com/fmiskovic/go-starter/pkg/validator"
	"github.com/matryer/is"
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
	assert := is.New(t)

	bunDb, app := test.SetUpServer(t)

	repo := NewRepo(bunDb)
	app.Delete("/user/:id", HandleDeleteById(repo))

	tests := []struct {
		name     string
		route    string
		given    func(t *testing.T) error
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:  "given user id should return 204 and delete user",
			route: "/user/1",
			given: func(t *testing.T) error {
				return repo.Create(context.Background(), &User{Email: "test1@fake.com"})
			},
			wantCode: 204,
			verify: func(t *testing.T, res *http.Response) {
				u, err := repo.GetById(context.Background(), 1)
				assert.True(err != nil)
				assert.True(u == nil)
			},
		},
		{
			name:  "given non-existing user id should return 204",
			route: "/user/11",
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 204,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
		{
			name:  "given zero user id should return 400",
			route: "/user/0",
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 400,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
		{
			name:  "given invalid user id should return 400",
			route: "/user/s",
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 400,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			err := tt.given(t)
			assert.NoErr(err)

			req := httptest.NewRequest("DELETE", tt.route, nil)

			// when
			res, err := app.Test(req, 1000)
			// then
			assert.NoErr(err)
			assert.Equal(res.StatusCode, tt.wantCode)
			tt.verify(t, res)
		})
	}
}

func TestHandleGetById(t *testing.T) {
	assert := is.New(t)

	bunDb, app := test.SetUpServer(t)

	repo := NewRepo(bunDb)
	app.Get("/user/:id", HandleGetById(repo))

	tests := []struct {
		name     string
		route    string
		given    func(t *testing.T) error
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:  "given user id should return 200 and user dto",
			route: "/user/1",
			given: func(t *testing.T) error {
				return repo.Create(context.Background(), &User{Email: "test1@fake.com", Gender: MALE})
			},
			wantCode: 200,
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
				assert.Equal(dto.Gender.numify(), MALE)
			},
		},
		{
			name:  "given non-existing user id should return 404",
			route: "/user/11",
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 404,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
		{
			name:  "given zero user id should return 400",
			route: "/user/0",
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 400,
			verify: func(t *testing.T, res *http.Response) {
			},
		},
		{
			name:  "given invalid user id should return 400",
			route: "/user/s",
			given: func(t *testing.T) error {
				return nil
			},
			wantCode: 400,
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
			name:  "given empty pageable should return 200",
			route: "/user",
			given: func(t *testing.T) error {
				return repo.Create(context.Background(), &User{Email: "test1@fake.com"})
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
				return repo.Create(context.Background(), &User{Email: "test2@fake.com"})
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
				return repo.Create(context.Background(), &User{Email: "test3@fake.com"})
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
