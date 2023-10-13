package model

import (
	"encoding/json"

	"github.com/mono83/maybe"
	maybejson "github.com/mono83/maybe/json"
)

// ResultPlace ...
type ResultPlace uint32

const (
	// InvalidPlace ...
	InvalidPlace ResultPlace = iota
	// FirstPlace ...
	FirstPlace
	// SecondPlace ...
	SecondPlace
	// ThrirdPlace ...
	ThrirdPlace
)

// String ...
func (p ResultPlace) String() string {
	switch p {
	case FirstPlace:
		return "🥇"
	case SecondPlace:
		return "🥈"
	case ThrirdPlace:
		return "🥉"
	}

	return ""
}

// GameResult ...
type GameResult struct {
	ID          int32
	GameID      int32
	ResultPlace ResultPlace
	RoundPoints maybe.Maybe[string]
}

// MarshalJSON ...
func (gr GameResult) MarshalJSON() ([]byte, error) {
	type wrapperGameResult struct {
		ID          int32
		GameID      int32
		ResultPlace ResultPlace
		RoundPoints maybejson.Maybe[string]
	}

	wgr := wrapperGameResult{
		ID:          gr.ID,
		GameID:      gr.GameID,
		ResultPlace: gr.ResultPlace,
		RoundPoints: maybejson.Wrap(gr.RoundPoints),
	}
	return json.Marshal(wgr)
}
