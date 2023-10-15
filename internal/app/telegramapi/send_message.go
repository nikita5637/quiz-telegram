package telegramapi

import (
	"context"

	telegrampb "github.com/nikita5637/quiz-telegram/pkg/pb/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendMessage ...
func (i *Implementation) SendMessage(ctx context.Context, req *telegrampb.SendMessageRequest) (*telegrampb.SendMessageResponse, error) {
	msg := tgbotapi.NewMessage(req.GetTelegramId(), req.GetMessage())
	_, err := i.bot.Send(msg)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return &telegrampb.SendMessageResponse{}, nil
}
