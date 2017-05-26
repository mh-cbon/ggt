package login

import "time"

// User data
type User struct {
	Login    string
	Password string
}

// GetID is useful for identity check.
func (t User) GetID() string {
	return t.Login
}

// HashedUser is an user with login capabilities
type HashedUser struct {
	User
	Hash       string
	LastLogin  time.Time
	LastLogout time.Time
}
