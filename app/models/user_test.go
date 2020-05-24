package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_GeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name    string
		u       *User
		wantErr bool
	}{
		{
			"valid password",
			&User{Password: "validpassword"},
			false,
		},
		{
			"empty password",
			&User{Password: ""},
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

func TestUser_CheckPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		u       *User
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

func GenerateTestUser(t *testing.T) *User {
	t.Helper()
	return &User{Name: "user", Password: "password"}
}

func TestUser_ValidateBeforeStoring(t *testing.T) {

	type fields struct {
		Name         string
		PasswordHash string
		Password     string
		Email        Email
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"valid",
			fields{"user", "hash", "password", Email{Address: "user@example.org", ActivationToken: "qerwwerhasdf"}},
			false,
		},
		{
			"empty activation token email",
			fields{"user", "hash", "password", Email{Address: "user@example", ActivationToken: ""}},
			true,
		},
		{
			"wrong email",
			fields{"user", "hash", "password", Email{Address: "user@example", ActivationToken: "qerwwerhasdf"}},
			true,
		},
		{
			"short pass",
			fields{"user", "hash", "sh", Email{Address: "user@example.org"}},
			true,
		},
		{
			"large pass",
			fields{"user", "hash", "superlargepasswordlol", Email{Address: "user@example.org"}},
			true,
		},
		{
			"no user",
			fields{"", "hash", "password", Email{Address: "user@example.org"}},
			true,
		},
		{
			"no hash",
			fields{"user", "", "password", Email{Address: "user@example.org"}},
			true,
		},
		{
			"no pass",
			fields{"user", "hash", "", Email{Address: "user@example.org"}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Name:         tt.fields.Name,
				PasswordHash: tt.fields.PasswordHash,
				Password:     tt.fields.Password,
				Email:        tt.fields.Email,
			}
			err := u.ValidateBeforeStoring()
			assert.Equal(t, (err != nil), tt.wantErr)
		})
	}
}
