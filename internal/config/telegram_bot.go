package config

import "fmt"

// TelegramBotConfig ...
type TelegramBotConfig struct {
	GamesWithPhotosListLimit byte   `toml:"games_with_photos_list_limit"`
	GroupID                  int64  `toml:"group_id"`
	RegistratorAPIAddress    string `toml:"registrator_api_address"`
	RegistratorAPIPort       uint16 `toml:"registrator_api_port"`
	TelegramAPIBindAddress   string `toml:"telegram_api_bind_address"`
	TelegramAPIBindPort      uint16 `toml:"telegram_api_bind_port"`
}

// GetTelegramAPIBindAddress ...
func GetTelegramAPIBindAddress() string {
	return fmt.Sprintf("%s:%d", globalConfig.TelegramAPIBindAddress, globalConfig.TelegramAPIBindPort)
}
