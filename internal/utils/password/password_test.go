package password

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		verify  func(password string, hash string, t *testing.T)
		wantErr bool
	}{
		{
			name: "given password should return hash",
			args: args{password: "Password1234!"},
			verify: func(password string, hash string, t *testing.T) {
				if !CheckPasswordHash(password, hash) {
					t.Errorf("Verify failed. The password: [%s] doesnot match the hash: [%s].", password, hash)
				}
			},
			wantErr: false,
		},
		{
			name:    "given long password should return error",
			args:    args{password: "Password1234!!!!Password1234!!!!Password1234!!!!Password1234!!!!Password1234!!!!"},
			verify:  func(password string, hash string, t *testing.T) {},
			wantErr: true,
		},
		{
			name: "given empty password should return hash",
			args: args{password: ""},
			verify: func(password string, hash string, t *testing.T) {
				if !CheckPasswordHash(password, hash) {
					t.Errorf("Verify failed. The password: [%s] doesnot match the hash: [%s].", password, hash)
				}
			},
			wantErr: false,
		},
		{
			name: "given blank password should return hash",
			args: args{password: " "},
			verify: func(password string, hash string, t *testing.T) {
				if !CheckPasswordHash(password, hash) {
					t.Errorf("Verify failed. The password: [%s] doesnot match the hash: [%s].", password, hash)
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.verify(tt.args.password, got, t)
		})
	}
}
