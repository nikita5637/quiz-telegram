package bot

// Command ...
type Command byte

// ADD NEW CONSTS AT THE BOTTOM ...
const (
	// CommandInvalid ...
	CommandInvalid Command = iota
	// CommandGamesList ...
	CommandGamesList
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
	// ...

	// CommandsNumber ...
	CommandsNumber
)

// TelegramRequest ...
type TelegramRequest struct {
	Command Command `json:"command,omitempty"`
	Body    []byte  `json:"body,omitempty"`
}

// Venue is data for CommandGetVenue
type VenueData struct {
	Title     string
	Address   string
	Latitude  float32
	Longitude float32
}
