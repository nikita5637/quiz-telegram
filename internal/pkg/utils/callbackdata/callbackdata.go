package callbackdata

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/nikita5637/quiz-telegram/internal/pkg/commands"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
)

// GetCallbackData ...
func GetCallbackData(ctx context.Context, command commands.Command, payload interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req := commands.TelegramRequest{
		Command: command,
		Body:    body,
	}

	callbackData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	if len(callbackData) > 64 {
		logger.ErrorKV(ctx, "callback data too long", "data", callbackData)
		return "", errors.New("callback data too long")
	}

	return string(callbackData), nil
}
