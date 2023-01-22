package model

import (
	"time"

	time_utils "github.com/nikita5637/quiz-registrator-api/utils/time"
	"github.com/nikita5637/quiz-telegram/internal/config"
)

// Game ...
type Game struct {
	ID          int32
	ExternalID  int32
	League      League
	Type        int32
	Number      string
	Name        string
	Place       Place
	Date        DateTime
	Price       uint32
	PaymentType string
	MaxPlayers  uint32
	Payment     PaymentType
	Registered  bool

	My                  bool
	NumberOfMyLegioners uint32
	NumberOfLegioners   uint32
	NumberOfPlayers     uint32
	ResultPlace         ResultPlace

	WithLottery bool
	DeletedAt   DateTime
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
