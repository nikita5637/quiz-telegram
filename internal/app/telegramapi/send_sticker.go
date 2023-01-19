package telegramapi

import (
	"context"

	telegrampb "github.com/nikita5637/quiz-telegram/pkg/pb/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendSticker ...
func (i *Implementation) SendSticker(ctx context.Context, req *telegrampb.SendStickerRequest) (*telegrampb.SendStickerResponse, error) {
	fileID := tgbotapi.FileID(req.GetStickerId())
	msg := tgbotapi.NewSticker(req.GetTelegramId(), fileID)
	_, err := i.bot.Send(msg)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return &telegrampb.SendStickerResponse{}, nil
}
