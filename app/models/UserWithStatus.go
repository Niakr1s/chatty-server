package models

// UserWithStatus ...
type UserWithStatus struct {
	User
	UserStatus
}

// NewUserWithStatus ...
func NewUserWithStatus(u User, s UserStatus) UserWithStatus {
	return UserWithStatus{User: u, UserStatus: s}
}
