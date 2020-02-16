package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	// ErrUserNilPass ...
	ErrUserNilPass = errors.New("(nil) passed instead of (User)")
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("(User) not found")
)

type Roles string

const (
	Admin          = "admin"
	NetworkManager = "manager"
	NetworkUser    = "user"
)

type User struct {
	gorm.Model        // Inject `ID`, `CreatedAt`, `UpdatedAt` and `DeletedAt`
	Username   string `json:"username" db:"username" sql:"not null"`
	Password   string `json:"password" db:"password" sql:"not null"`
	Email      string `json:"email" db:"email"`
	Role       Roles  `json:"role" sql:"type:ENUM('admin', 'manager', 'user')" gorm:"default:'user'"`
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
