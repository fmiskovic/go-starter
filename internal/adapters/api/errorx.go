package api

import "errors"

var (
	ErrParseReqBody      = errors.New("failed to parse request body")
	ErrEntityCreate      = errors.New("failed to create entity")
	ErrEntityUpdate      = errors.New("failed to update entity")
	ErrGetById           = errors.New("failed to get entity by id")
	ErrInvalidId         = errors.New("invalid id")
	ErrDeleteById        = errors.New("failed to delete entity by id")
	ErrInvalidPageSize   = errors.New("invalid page size number")
	ErrInvalidPageOffset = errors.New("invalid page offset number")
	ErrGetPage           = errors.New("failed to get entities page")
)

// ErrorX represents a custom error struct that contains optionally service and application error.
type ErrorX struct {
	srvErr error
	appErr error
}

// Error is implementation of error interface.
func (x ErrorX) Error() string {
	if x.srvErr == nil && x.appErr == nil {
		return ""
	}
	return errors.Join(x.srvErr, x.appErr).Error()
}

// New instantiate new ErrorX.
func New(opts ...Option) *ErrorX {
	x := &ErrorX{}
	for _, opt := range opts {
		opt(x)
	}
	return x
}

// Option func used to construct ErrorX.
type Option func(x *ErrorX)

func WithSvcErr(svcErr error) Option {
	return func(x *ErrorX) {
		x.srvErr = svcErr
	}
}

func WithAppErr(appErr error) Option {
	return func(x *ErrorX) {
		x.appErr = appErr
	}
}
