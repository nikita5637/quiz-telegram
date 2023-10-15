package config

import "github.com/spf13/viper"

func initBotConfigureParams() {
	_ = viper.BindEnv("bot.ics_manager_api.address")
	_ = viper.BindEnv("bot.registrator_api.address")
	_ = viper.BindEnv("bot.token")
}
