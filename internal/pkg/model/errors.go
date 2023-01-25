package model

import "errors"

var (
	// ErrGameNotFound ...
	ErrGameNotFound = errors.New("game not found")
	// ErrLeagueNotFound ...
	ErrLeagueNotFound = errors.New("league not found")
	// ErrPlaceNotFound ...
	ErrPlaceNotFound = errors.New("place not found")
)
