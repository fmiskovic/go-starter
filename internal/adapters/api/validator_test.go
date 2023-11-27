package api

import (
	"reflect"
	"testing"
	"time"
)

type TestData struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Email       string    `validate:"required,min=3" json:"email"`
	FullName    string    `json:"fullname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      GenderDto `json:"gender"`
	Enabled     bool      `json:"enabled"`
}

func TestValidator_Validate(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "given valid data should return no errors",
			args: args{data: TestData{Email: "test@fake.com"}},
			want: []string{},
		},
		{
			name: "given invalid email should return error",
			args: args{data: TestData{Email: "t@"}},
			want: []string{"[Email]: 't@' | Needs to implement 'min'"},
		},
		{
			name: "given nil email should return error",
			args: args{data: TestData{}},
			want: []string{"[Email]: '' | Needs to implement 'required'"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			if got := v.Validate(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
