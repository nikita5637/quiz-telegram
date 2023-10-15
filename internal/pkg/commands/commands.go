package commands

import (
	"encoding/json"

	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// Command ...
type Command byte

// ADD NEW CONSTS AT THE BOTTOM ...
const (
	// CommandInvalid ...
	CommandInvalid Command = iota
	// CommandGetGamesList ...
	CommandGetGamesList
	// CommandGetGame ...
	CommandGetGame
	// CommandRegisterGame ...
	CommandRegisterGame
	// CommandRegisterPlayer ...
	CommandRegisterPlayer
	// CommandUnregisterGame ...
	CommandUnregisterGame
	// CommandUnregisterPlayer ...
	CommandUnregisterPlayer
	// CommandPlayersListByGame ...
	CommandPlayersListByGame
	// CommandGetGamePhotos ...
	CommandGetGamePhotos
	_deleted_CommandGetPassedAndRegisteredGamesListNextPage // nolint:revive
	_deleted_CommandGetPassedAndRegisteredGamesListPrevPage // nolint:revive
	// CommandUpdatePayment ...
	CommandUpdatePayment
	// CommandChangeEmail ...
	CommandChangeEmail
	// CommandChangeName ...
	CommandChangeName
	// CommandChangePhone ...
	CommandChangePhone
	// CommandLottery ...
	CommandLottery
	// CommandGetVenue ...
	CommandGetVenue
	// CommandChangeBirthdate ...
	CommandChangeBirthdate
	// CommandChangeSex ...
	CommandChangeSex
	// CommandUpdatePlayerRegistration ...
	CommandUpdatePlayerRegistration
	// CommandGetPassedAndRegisteredGames for pagination
	CommandGetPassedAndRegisteredGamesList
	// CommandGetRegisteredGamesList ...
	CommandGetRegisteredGamesList
	// CommandGetUserGamesList ...
	CommandGetUserGamesList
	// ...

	// CommandsNumber ...
	CommandsNumber
)

// TelegramRequest ...
type TelegramRequest struct {
	Command Command         `json:"c,omitempty"`
	Body    json.RawMessage `json:"b,omitempty"`
}

// GetGameData is a payload for command CommandGetGame
type GetGameData struct {
	GameID                  int32   `json:"g,omitempty"`
	PageIndex               uint32  `json:"p,omitempty"`
	GetRootGamesListCommand Command `json:"c,omitempty"`
}

// GetGamePhotosData is a payload for command CommandGetGamePhotos
type GetGamePhotosData struct {
	GameID int32 `json:"g,omitempty"`
}

// GetGamesListData is a payload for command CommandGetGamesList
type GetGamesListData struct {
	Command Command `json:"c,omitempty"`
}

// GetPassedAndRegisteredGamesListData is a payload for commands CommandGetPassedAndRegisteredGamesListNextPage and CommandGetPassedAndRegisteredGamesListPrevPage
type GetPassedAndRegisteredGamesListData struct {
	Page     uint64 `json:"p,omitempty"`
	PageSize uint64 `json:"ps,omitempty"`
}

// GetVenueData is a payload for command CommandGetVenue
type GetVenueData struct {
	PlaceID int32 `json:"p,omitempty"`
}

// LotteryData is a payload for command CommandLottery
type LotteryData struct {
	GameID                  int32   `json:"g,omitempty"`
	GetRootGamesListCommand Command `json:"c,omitempty"`
}

// PlayersListByGameData is a payload for command CommandPlayersListByGame
type PlayersListByGameData struct {
	GameID int32 `json:"g,omitempty"`
}

// RegisterGameData is a payload for command CommandRegisterGame
type RegisterGameData struct {
	GameID                  int32   `json:"g,omitempty"`
	GetRootGamesListCommand Command `json:"c,omitempty"`
}

// RegisterPlayerData is a payload for command CommandRegisterPlayer
type RegisterPlayerData struct {
	GameID                  int32        `json:"g,omitempty"`
	UserID                  int32        `json:"u,omitempty"`
	RegisteredBy            int32        `json:"r,omitempty"`
	Degree                  model.Degree `json:"d,omitempty"`
	GetRootGamesListCommand Command      `json:"c,omitempty"`
}

// UnregisterGameData is a payload for command CommandUnregisterGame
type UnregisterGameData struct {
	GameID                  int32   `json:"g,omitempty"`
	GetRootGamesListCommand Command `json:"c,omitempty"`
}

// UnregisterPlayerData is a payload for command CommandUnregisterPlayer
type UnregisterPlayerData struct {
	GameID                  int32        `json:"g,omitempty"`
	UserID                  int32        `json:"u,omitempty"`
	RegisteredBy            int32        `json:"r,omitempty"`
	Degree                  model.Degree `json:"d,omitempty"`
	GetRootGamesListCommand Command      `json:"c,omitempty"`
}

// UpdatePaymentData is a payload for command CommandUpdatePayment
type UpdatePaymentData struct {
	GameID                  int32   `json:"g,omitempty"`
	Payment                 int32   `json:"p,omitempty"`
	GetRootGamesListCommand Command `json:"c,omitempty"`
}

// UpdatePlayerRegistration is a payload for command CommandUpdatePlayerRegistration
type UpdatePlayerRegistration struct {
	GameID                  int32        `json:"g,omitempty"`
	UserID                  int32        `json:"u,omitempty"`
	RegisteredBy            int32        `json:"r,omitempty"`
	Degree                  model.Degree `json:"d,omitempty"`
	GetRootGamesListCommand Command      `json:"c,omitempty"`
}
