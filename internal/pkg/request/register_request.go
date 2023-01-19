package request

import (
	"context"
	"fmt"

	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

// RegisterRequest ...
func (f *Facade) RegisterRequest(ctx context.Context, request model.Request) (string, error) {
	requestUUID := ""
	if _, err := f.requestStorage.GetRequestByUUID(ctx, request.UUID); err != nil {
		if err == model.ErrRequestNotFound {
			_, requestUUID, err = f.requestStorage.Insert(ctx, request)
			if err != nil {
				return "", fmt.Errorf("register request error: %w", err)
			}
		} else {
			return "", fmt.Errorf("register request error: %w", err)
		}
	}

	f.cache[request.UUID] = request

	return requestUUID, nil
}
