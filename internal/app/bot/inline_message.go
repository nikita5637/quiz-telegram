package bot

import (
	"context"
	"strconv"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	youAreAlreadyRegisteredForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_already_registered_for_the_game",
		FallBack: "You are already registered for the game",
	}
	youAreRegisteredForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_registered_for_the_game",
		FallBack: "You are registered for the game",
	}
)

// HandleInlineMessage ...
func (b *Bot) HandleInlineMessage(ctx context.Context, update *tgbotapi.Update) error {
	clientID := update.CallbackQuery.From.ID

	logger.DebugKV(ctx, "new inline message incoming", "clientID", clientID, "data", update.CallbackData())

	ctx = telegram_utils.NewContextWithClientID(ctx, clientID)

	return b.handleInlineMessage(ctx, update)
}

func (b *Bot) handleInlineMessage(ctx context.Context, update *tgbotapi.Update) error {
	gameID, err := strconv.ParseUint(update.CallbackQuery.Data, 10, 32)
	if err != nil {
		return err
	}

	req := &registrator.RegisterPlayerRequest{
		GameId:     int32(gameID),
		PlayerType: registrator.PlayerType_PLAYER_TYPE_MAIN,
		Degree:     registrator.Degree_DEGREE_LIKELY,
	}

	resp, err := b.registratorServiceClient.RegisterPlayer(ctx, req)
	if err != nil {
		st := status.Convert(err)

		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					cb := tgbotapi.NewCallback(update.CallbackQuery.ID, localizedMessage)
					_, err = b.bot.Request(cb)
					return err
				}
			}
		} else if st.Code() == codes.AlreadyExists {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					cb := tgbotapi.NewCallback(update.CallbackQuery.ID, localizedMessage)
					_, err = b.bot.Request(cb)
					return err
				}
			}
		}

		return err
	}

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, resp.GetStatus().String())
	switch resp.GetStatus() {
	case registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_ALREADY_REGISTERED:
		cb = tgbotapi.NewCallback(update.CallbackQuery.ID, getTranslator(youAreAlreadyRegisteredForTheGameLexeme)(ctx))
	case registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_OK:
		cb = tgbotapi.NewCallback(update.CallbackQuery.ID, getTranslator(youAreRegisteredForTheGameLexeme)(ctx))
	}

	_, err = b.bot.Request(cb)
	return err
}
