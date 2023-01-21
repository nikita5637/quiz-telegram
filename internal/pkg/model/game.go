package model

import (
	"time"

	time_utils "github.com/nikita5637/quiz-registrator-api/utils/time"
	"github.com/nikita5637/quiz-telegram/internal/config"
)

// Game ...
type Game struct {
	Date       DateTime
	ExternalID int32
	GameType   int32
	ID         int32
	LeagueID   int32
	MaxPlayers byte
	Number     string
	PlaceID    int32
	Place      Place
	Registered bool
	Payment    PaymentType
	gameAdditionalInfo
}

type gameAdditionalInfo struct {
	DeletedAt       DateTime    `json:"deleted_at,omitempty"`
	My              bool        `json:"my,omitempty"`
	MyLegioners     byte        `json:"my_legioners,omitempty"`
	NumberLegioners byte        `json:"number_legioners,omitempty"`
	NumberPlayers   byte        `json:"number_players,omitempty"`
	ResultPlace     ResultPlace `json:"result_place,omitempty"`
	WithLottery     bool
}

// DateTime ...
func (g *Game) DateTime() DateTime {
	if g == nil {
		return DateTime{}
	}

	return g.Date
}

// IsActive ...
func (g *Game) IsActive() bool {
	if g == nil {
		return false
	}

	activeGameLag := config.GetValue("ActiveGameLag").Uint16()
	return g.DeletedAt.AsTime().IsZero() && time_utils.TimeNow().Before(g.DateTime().AsTime().Add(time.Duration(activeGameLag)*time.Second))
}
