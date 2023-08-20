package gameplayers

import "errors"

var (
	// ErrGamePlayerAlreadyRegistered ...
	ErrGamePlayerAlreadyRegistered = errors.New("game player already registered")
	// ErrGamePlayerNotFound ...
	ErrGamePlayerNotFound = errors.New("game player not found")
	// ErrNoFreeSlot ...
	ErrNoFreeSlot = errors.New("no free slot")

	// ReasonNoFreeSlot ...
	ReasonNoFreeSlot = "THERE_ARE_NO_FREE_SLOT"
)
