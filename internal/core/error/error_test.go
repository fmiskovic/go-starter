package error

import (
	"errors"
	"testing"
)

func TestApiError(t *testing.T) {
	type fields struct {
		srvErr error
		appErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "given server and app error should return joined error message",
			fields: fields{srvErr: errors.New("svcErr"), appErr: errors.New("appErr")},
			want:   errors.Join(errors.New("svcErr"), errors.New("appErr")).Error(),
		},
		{
			name:   "given server error only should return server error message",
			fields: fields{srvErr: errors.New("svcErr"), appErr: nil},
			want:   errors.New("svcErr").Error(),
		},
		{
			name:   "given app error only should return app error message",
			fields: fields{srvErr: nil, appErr: errors.New("appErr")},
			want:   errors.New("appErr").Error(),
		},
		{
			name:   "given none error only should return empty error message",
			fields: fields{srvErr: nil, appErr: nil},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := New(WithSvcErr(tt.fields.srvErr), WithAppErr(tt.fields.appErr))
			if got := x.Error(); got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
