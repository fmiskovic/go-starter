package ports

import (
	"context"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
)

type UserService[ID any] interface {
	SingIn(ctx context.Context, req *user.SignInRequest) (*user.SignInResponse, error)
	SingUp(ctx context.Context, req *user.CreateRequest) (*user.SignUpResponse, error)
	ConfirmEmail(ctx context.Context, req user.ConfirmEmailRequest) error
	Create(ctx context.Context, req *user.CreateRequest) (*user.CreateResponse, error)
	Update(ctx context.Context, req *user.UpdateRequest) (*user.UpdateResponse, error)
	GetById(ctx context.Context, id ID) (*user.Dto, error)
	DeleteById(ctx context.Context, id ID) error
	GetPage(ctx context.Context, pagabale domain.Pageable) (*domain.Page[user.Dto], error)
	AddRoles(ctx context.Context, roles []string, id ID) error
	RemoveRoles(ctx context.Context, roles []string, id ID) error
	EnableDisable(ctx context.Context, id ID) error
	ChangePassword(ctx context.Context, req *user.ChangePasswordRequest) error
}
