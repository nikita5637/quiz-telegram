package model

import (
	"github.com/mono83/maybe"
)

// Degree ...
type Degree int32

const (
	// DegreeInvalid ...
	DegreeInvalid Degree = iota
	// DegreeLikely ...
	DegreeLikely
	// DegreeUnlikely ...
	DegreeUnlikely
)

// GamePlayer ...
type GamePlayer struct {
	ID           int32
	GameID       int32
	UserID       maybe.Maybe[int32]
	RegisteredBy int32
	Degree       Degree
}
