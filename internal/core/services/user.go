package services

import (
	"context"
	"errors"
	"time"

	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/security"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/fmiskovic/go-starter/internal/utils/password"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserService.
type UserService struct {
	repo       ports.UserRepo[uuid.UUID]
	authConfig configs.AuthConfig
}

// NewUserService instantiate new UserService.
func NewUserService(userRepo ports.UserRepo[uuid.UUID], authConfig configs.AuthConfig) UserService {
	return UserService{userRepo, authConfig}
}

// SingIn authenticates user.
// Returns new signed jwt token.
func (s UserService) SingIn(ctx context.Context, req user.SignInRequest) (*user.SignInResponse, error) {
	u, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !password.CheckPasswordHash(req.Password, u.Credentials.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"email": u.Email,
		"sub":   u.ID,
		"name":  u.FullName,
		"roles": u.Roles,
		"exp":   time.Now().Add(time.Hour * s.authConfig.TokenExp).Unix(),
		"iat":   time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Generate signed token and send it as response.
	signedToken, err := token.SignedString([]byte(s.authConfig.Secret))
	if err != nil {
		return nil, err
	}

	return &user.SignInResponse{Token: signedToken}, nil
}

// ConfirmEmail enables user when user confirs it's email address.
func (s UserService) ConfirmEmail(ctx context.Context, req user.ConfirmEmailRequest) error {
	return nil
}

// SingUp register new user.
func (s UserService) SingUp(ctx context.Context, req user.CreateRequest) (*user.SignUpResponse, error) {
	u, err := s.createUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &user.SignUpResponse{ID: u.ID.String()}, nil
}

// Create creates new user.
// This function is for admin user only.
// Returns newly created user.
func (s UserService) Create(ctx context.Context, req user.CreateRequest) (*user.CreateResponse, error) {
	u, err := s.createUser(ctx, req)
	if err != nil {
		return nil, err
	}

	u, err = s.repo.GetById(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	return &user.CreateResponse{Dto: *user.ConvertToDto(u)}, nil
}

// Update updates existing user.
// Returns user with fresh changes.
func (s UserService) Update(ctx context.Context, req user.UpdateRequest) (*user.UpdateResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}
	u := user.New(
		user.Id(id),
		user.Email(req.Email),
		user.DateOfBirth(req.DateOfBirth),
		user.FullName(req.FullName),
		user.Location(req.Location),
		user.Sex(req.Gender.Numberfy()),
	)
	if err = s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	u, err = s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user.UpdateResponse{Dto: *user.ConvertToDto(u)}, nil
}

// GetById returns existing user.
func (s UserService) GetById(ctx context.Context, id uuid.UUID) (*user.Dto, error) {
	u, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ConvertToDto(u), nil
}

// DeleteById deletes existing user.
func (s UserService) DeleteById(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteById(ctx, id)
}

// GetPage returns page of users.
func (s UserService) GetPage(ctx context.Context, pagabale domain.Pageable) (*domain.Page[user.Dto], error) {
	page, err := s.repo.GetPage(ctx, pagabale)
	if err != nil {
		return nil, err
	}
	return user.ConvertToPageDto(page), nil
}

// AddRoles appends user roles.
func (s UserService) AddRoles(ctx context.Context, roles []string, id uuid.UUID) error {
	u, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	// TODO: implement
	for _, name := range roles {
		role := security.NewRole(name)
		role.UserID = u.ID
		u.Roles = append(u.Roles, role)
	}

	return nil
}

// ChangePassword updates user password.
func (s UserService) ChangePassword(ctx context.Context, req user.ChangePasswordRequest) error {
	// TODO: implement
	return nil
}

// EnableDisable is for admin usage only, to enable user if disabled and vice versa.
func (s UserService) EnableDisable(ctx context.Context, id uuid.UUID) error {
	// TODO: implement
	return nil
}

func (s UserService) createUser(ctx context.Context, req user.CreateRequest) (*user.User, error) {
	crd := security.NewCredentials(req.Username, req.Password)
	role := security.NewRole(security.ROLE_USER)
	u := user.New(
		user.Email(req.Email),
		user.DateOfBirth(req.DateOfBirth),
		user.FullName(req.FullName),
		user.Location(req.Location),
		user.Sex(req.Gender.Numberfy()),
		user.Credentials(crd),
		user.Roles(role),
	)

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
