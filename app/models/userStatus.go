package models

// UserStatus ...
type UserStatus struct {
	Verified bool `json:"verified"`
	Admin    bool `json:"admin"`
}
