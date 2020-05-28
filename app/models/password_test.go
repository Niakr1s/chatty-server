package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generatePasswordHash(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid password",
			args{"validpassword"},
			false,
		},
		{
			"empty password",
			args{""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generatePasswordHash(tt.args.pass)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func GenerateTestUser(t *testing.T) FullUser {
	t.Helper()
	return NewFullUser("user", "user@example.com", "password")
}

func TestUser_CheckPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		u       FullUser
		args    args
		wantErr bool
	}{
		{
			"same password",
			GenerateTestUser(t),
			args{"password"},
			false,
		},
		{
			"other password",
			GenerateTestUser(t),
			args{"otherpassword"},
			true,
		},
		{
			"empty password",
			GenerateTestUser(t),
			args{""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.GeneratePasswordHash()
			err := tt.u.CheckPassword(tt.args.password)
			assert.Equal(t, err != nil, tt.wantErr, err)
		})
	}
}

func TestUser_GeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name    string
		u       *FullUser
		wantErr bool
	}{
		{
			"valid password",
			&FullUser{Pass: Pass{Password: "validpassword"}},
			false,
		},
		{
			"empty password",
			&FullUser{Pass: Pass{Password: ""}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.u.GeneratePasswordHash()

			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
