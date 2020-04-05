package models

import (
	"errors"
)

var (
	// ErrNetworkNilPass ...
	ErrNetworkNilPass = errors.New("(nil) passed instead of (Network)")
	// ErrNetworkNotFound ...
	ErrNetworkNotFound = errors.New("(Network) not found")
)

// Network will be a bridge struct between a User and Nodes in the current
// Network.
//  - `Network` can have zero or more nodes. tl;dr Has Many `Nodes`.
//  - `Network` can belong to only one user. tl;dr Belongs To `User`.
//  - `Network` can be accessed by many users. Not implemented yet.
type Network struct {
	ID   uint   `gorm:"primary_key"`
	Name string `json:"name"`

	// Refers to UserId to whom it belongs to.
	UserRefer uint `json:"user_refer"`

	// a `Has Many` relationship. Network 1 <-> 0..* Node.
	// Network `Has Many` Nodes.
	Nodes []Node `json:"nodes,omitempty" gorm:"foreignkey:NetworkRefer"`

	// a `Has Many` relationship. Network 1 <-> 0..* UnregisteredNode.
	// Network `Has Many` UnregisteredNodes.
	UnregisteredNodes []UnregisteredNode `json:"unregistered_nodes,omitempty" gorm:"foreignkey:NetworkRefer"`
}
