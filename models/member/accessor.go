package member

import (
	"errors"
)

// ErrMemberNotExist occurs when the member with the member.ID was not found.
var ErrMemberNotExist = errors.New("member does not exist. check member.ID")

// Accessor is an interface to access the tasks.
type Accessor interface {
	Get(id ID) (Member, error)
	// Put(id ID, m Member) error
	// Post(id ID) (NameBirth, error)
	// Delete(id ID) error

	//
	GetExistIDs() []ID
}
