package user

import (
	"bytes"
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
		reqBody  *Request
		wantCode int
		verify   func(t *testing.T, res *http.Response)
	}{
		{
			name:     "given valid user request should return 201",
			route:    "/user",
			reqBody:  NewRequest(Email("test1@fake.com")),
			wantCode: 201,
			verify: func(t *testing.T, res *http.Response) {
				resBody := res.Body
				defer func(body io.ReadCloser) {
					if err := body.Close(); err != nil {
						fmt.Println("error occurred on body close:", err.Error())
					}
				}(resBody)

				u := &User{}
				err := json.NewDecoder(resBody).Decode(u)
				assert.NoErr(err)
				assert.Equal(u.Email, "test1@fake.com")
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
			reqBody:  NewRequest(Email("")),
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
			//if got := HandleGetPage(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("HandleGetPage() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestHandleUpdate(t *testing.T) {
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
			//if got := HandleUpdate(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("HandleUpdate() = %v, want %v", got, tt.want)
			//}
		})
	}
}
