package team

import (
	"errors"
)

// ErrTeamNotExist occurs when the team with the member.ID was not found.
var ErrTeamNotExist = errors.New("team does not exist. check team.ID")

// Accessor is an interface to access the tasks.
type Accessor interface {
	Get(id ID) (Team, error)
	// Put(id ID, m Member) error
	// Post(id ID) (NameBirth, error)
	// Delete(id ID) error

	//
	GetExistIDs() []ID
}
