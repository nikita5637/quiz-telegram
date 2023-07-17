package commands

// Command ...
type Command byte

// ADD NEW CONSTS AT THE BOTTOM ...
const (
	// CommandInvalid ...
	CommandInvalid Command = iota
	// CommandGetGamesList ...
	CommandGetGamesList // TODO delete
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
	// CommandGetListGamesWithPhotosNextPage ...
	CommandGetListGamesWithPhotosNextPage
	// CommandGetListGamesWithPhotosPrevPage ...
	CommandGetListGamesWithPhotosPrevPage
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
	// ...

	// CommandsNumber ...
	CommandsNumber
)

// TelegramRequest ...
type TelegramRequest struct {
	Command Command `json:"c,omitempty"`
	Body    []byte  `json:"b,omitempty"`
}

// GetGameData is a payload for command CommandGetGame
type GetGameData struct {
	GameID    int32  `json:"g,omitempty"`
	PageIndex uint32 `json:"p,omitempty"`
}

// GetGamePhotosData is a payload for command CommandGetGamePhotos
type GetGamePhotosData struct {
	GameID int32 `json:"g,omitempty"`
}

// GetGamesListData is a payload for command CommandGetGamesList
type GetGamesListData struct {
	Active bool `json:"a,omitempty"`
}

// GetGamesWithPhotosData is a payload for commands CommandGetListGamesWithPhotosNextPage and CommandGetListGamesWithPhotosPrevPage
type GetGamesWithPhotosData struct {
	Limit  uint32 `json:"l,omitempty"`
	Offset uint32 `json:"o,omitempty"`
}

// GetVenueData is a payload for command CommandGetVenue
type GetVenueData struct {
	PlaceID int32 `json:"p,omitempty"`
}

// LotteryData is a payload for command CommandLottery
type LotteryData struct {
	GameID int32 `json:"g,omitempty"`
}

// PlayersListByGameData is a payload for command CommandPlayersListByGame
type PlayersListByGameData struct {
	GameID int32 `json:"g,omitempty"`
}

// RegisterGameData is a payload for command CommandRegisterGame
type RegisterGameData struct {
	GameID int32 `json:"g,omitempty"`
}

// RegisterPlayerData is a payload for command CommandRegisterPlayer
type RegisterPlayerData struct {
	GameID     int32 `json:"g,omitempty"`
	PlayerType int32 `json:"p,omitempty"`
	Degree     int32 `json:"d,omitempty"`
}

// UnregisterGameData is a payload for command CommandUnregisterGame
type UnregisterGameData struct {
	GameID int32 `json:"g,omitempty"`
}

// UnregisterPlayerData is a payload for command CommandUnregisterPlayer
type UnregisterPlayerData struct {
	GameID     int32 `json:"g,omitempty"`
	PlayerType int32 `json:"p,omitempty"`
}

// UpdatePaymentData is a payload for command CommandUpdatePayment
type UpdatePaymentData struct {
	GameID  int32 `json:"g,omitempty"`
	Payment int32 `json:"p,omitempty"`
}
