package api

import "errors"

var (
	NilUserError         = errors.New("user entity can not be nil")
	ErrUserReqBody       = errors.New("failed to parse user request body")
	ErrUserCreate        = errors.New("failed to create user")
	ErrUserUpdate        = errors.New("failed to update user")
	ErrUserGetById       = errors.New("failed to get user by id")
	ErrInvalidUserId     = errors.New("invalid user id")
	ErrUserDeleteById    = errors.New("failed to delete user by id")
	ErrInvalidPageSize   = errors.New("invalid page size number")
	ErrInvalidPageOffset = errors.New("invalid page offset number")
	ErrUserGetPage       = errors.New("failed to get users page")
)
