package model

import (
	"encoding/json"

	"github.com/mono83/maybe"
	maybejson "github.com/mono83/maybe/json"
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

// MarshalJSON ...
func (gp GamePlayer) MarshalJSON() ([]byte, error) {
	type wrapperGamePlayer struct {
		ID           int32
		GameID       int32
		UserID       maybejson.Maybe[int32]
		RegisteredBy int32
		Degree       Degree
	}

	wgp := wrapperGamePlayer{
		ID:           gp.ID,
		GameID:       gp.GameID,
		UserID:       maybejson.Wrap(gp.UserID),
		RegisteredBy: gp.RegisteredBy,
		Degree:       gp.Degree,
	}
	return json.Marshal(wgp)
}
