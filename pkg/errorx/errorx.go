package errorx

import "errors"

type ErrorX struct {
	srvErr error
	appErr error
}

func (x ErrorX) Error() string {
	if x.srvErr == nil && x.appErr == nil {
		return ""
	}
	return errors.Join(x.srvErr, x.appErr).Error()
}

func New(opts ...Option) *ErrorX {
	x := &ErrorX{}
	for _, opt := range opts {
		opt(x)
	}
	return x
}

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
