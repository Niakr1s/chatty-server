package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			u := &FullUser{
				User: User{UserName: tt.fields.Name},
				Pass: Pass{
					PasswordHash: tt.fields.PasswordHash,
					Password:     tt.fields.Password,
				},
				Email: tt.fields.Email,
			}
			err := u.ValidateBeforeStoring()
			assert.Equal(t, (err != nil), tt.wantErr)
		})
	}
}
