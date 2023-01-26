package bot

import (
	"context"
	"errors"
	"strconv"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	noFreeSlotLexeme = i18n.Lexeme{
		Key:      "no_free_slot",
		FallBack: "There are not free slot",
	}
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

	registerStatus, err := b.gamesFacade.RegisterPlayer(ctx, int32(gameID), int32(registrator.PlayerType_PLAYER_TYPE_MAIN), int32(registrator.Degree_DEGREE_LIKELY))
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			cb := tgbotapi.NewCallback(update.CallbackQuery.ID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Request(cb)
			return err
		} else if errors.Is(err, model.ErrNoFreeSlot) {
			cb := tgbotapi.NewCallback(update.CallbackQuery.ID, getTranslator(noFreeSlotLexeme)(ctx))
			_, err = b.bot.Request(cb)
			return err
		}

		return err
	}

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	switch registrator.RegisterPlayerStatus(registerStatus) {
	case registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_ALREADY_REGISTERED:
		cb = tgbotapi.NewCallback(update.CallbackQuery.ID, getTranslator(youAreAlreadyRegisteredForTheGameLexeme)(ctx))
	case registrator.RegisterPlayerStatus_REGISTER_PLAYER_STATUS_OK:
		cb = tgbotapi.NewCallback(update.CallbackQuery.ID, getTranslator(youAreRegisteredForTheGameLexeme)(ctx))
	}

	_, err = b.bot.Request(cb)
	return err
}
