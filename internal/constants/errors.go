package constants

import (
	"errors"
)

var Errs *Errors

func init() {
	Errs = NewErrors()
}

type Errors struct {
	Unknown,
	InvalidArgument,
	MatchAlreadyStarted,
	NoVChUser,
	NoAvailableVCh,
	MatchNotFound,
	ConflictVCh,
	InvalidTeam error
}

func NewErrors() *Errors {
	errs := &Errors{
		Unknown:             errors.New("unknown error"),
		InvalidArgument:     errors.New("invalid Argument"),
		MatchAlreadyStarted: errors.New("match has already started"),
		NoVChUser:           errors.New("no voice channel user"),
		NoAvailableVCh:      errors.New("no available voice channel"),
		MatchNotFound:       errors.New("match not found"),
		ConflictVCh:         errors.New("conflict voice channel usage"),
		InvalidTeam:         errors.New("invalid team"),
	}
	return errs
}
