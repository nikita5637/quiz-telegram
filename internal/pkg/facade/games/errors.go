package games

import "errors"

var (
	// ErrGameHasPassed ...
	ErrGameHasPassed = errors.New("game has passed")
	// ErrGameNotFound ...
	ErrGameNotFound = errors.New("game not found")

	// ReasonGameHasPassed ...
	ReasonGameHasPassed = "GAME_HAS_PASSED"
	// ReasonGameNotFound ...
	ReasonGameNotFound = "GAME_NOT_FOUND"
)
