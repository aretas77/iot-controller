package models

import "errors"

var (
	// ErrUserNilPass ...
	ErrUserNilPass = errors.New("(nil) passed instead of (User)")
	// ErrUserNotFound ...
	ErrNodeNotFound = errors.New("(User) not found")
)

type User struct {
	ID       string `json:"_key" db:"_key"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
}

type UserService interface {

	// Init should initialize all required information and migration data.
	Init() error

	// Create should create a User and return a created User.
	Create(userID string, u *User) (*User, error)

	// Get should return a User by given ID.
	Get(userID string) (*User, error)

	// All should return all Users.
	All() ([]User, error)

	// Update should update the given User with given new values.
	Update(u *User) (*User, error)
}
