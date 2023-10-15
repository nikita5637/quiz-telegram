package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	// amqp://guest:guest@localhost:5672/
	rabbitMQURLFormat = "amqp://%s:%s@%s:%d/"
)

func initRemindManagerConfigureParams() {
	_ = viper.BindEnv("rabbitmq.address")
	_ = viper.BindEnv("rabbitmq.credentials.password")
}

// GetRabbitMQURL ...
func GetRabbitMQURL() string {
	return fmt.Sprintf(rabbitMQURLFormat,
		viper.GetString("rabbitmq.credentials.username"),
		viper.GetString("rabbitmq.credentials.password"),
		viper.GetString("rabbitmq.address"),
		viper.GetUint32("rabbitmq.port"),
	)
}
