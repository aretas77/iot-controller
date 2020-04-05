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
	// ErrUserUnauthorized ...
	ErrUserUnauthorized = errors.New("(User) not authorized")
)

const (
	Admin          = "admin"
	NetworkManager = "manager"
	NetworkUser    = "user"
)

// Credentials ...
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User ...
type User struct {
	gorm.Model        // Inject `ID`, `CreatedAt`, `UpdatedAt` and `DeletedAt`
	Username   string `json:"username" gorm:"username" sql:"not null"`
	Password   string `json:"password" sql:"not null"`
	Email      string `json:"email"`
	Role       string `json:"role" sql:"type:ENUM('admin', 'manager', 'user')" gorm:"default:'user'"`

	// XXX: Should be kept in some cache storage, e.g. Redis.
	Token string `json:"token"`

	// a `Has Many` relationship. User 1 <-> 0..* Network.
	// User `Has Many` Networks.
	Networks []Network `json:"networks" gorm:"foreignkey:UserRefer"`
}
