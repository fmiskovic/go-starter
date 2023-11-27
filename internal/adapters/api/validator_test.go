package api

import (
	"reflect"
	"testing"
)

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
			args: args{data: NewUserDto(Email("test@fake.com"))},
			want: []string{},
		},
		{
			name: "given invalid email should return error",
			args: args{data: NewUserDto(Email("t@"))},
			want: []string{"[Email]: 't@' | Needs to implement 'min'"},
		},
		{
			name: "given nil email should return error",
			args: args{data: NewUserDto()},
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
