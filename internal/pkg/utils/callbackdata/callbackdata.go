package callbackdata

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nikita5637/quiz-telegram/internal/pkg/commands"
)

// GetCallbackData ...
func GetCallbackData(ctx context.Context, command commands.Command, payload interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshaling payload error: %w", err)
	}

	req := commands.TelegramRequest{
		Command: command,
		Body:    body,
	}

	callbackData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshaling callback data error: %w", err)
	}

	if len(callbackData) > 64 {
		return "", fmt.Errorf("callback data len is too long: %s", callbackData)
	}

	return string(callbackData), nil
}
