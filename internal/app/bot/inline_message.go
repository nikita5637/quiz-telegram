package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mono83/maybe"
	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gameplayers"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"
	userutils "github.com/nikita5637/quiz-telegram/utils/user"

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
	ctx = telegram_utils.NewContextWithClientID(ctx, clientID)

	logger.DebugKV(ctx, "new inline message incoming", "clientID", clientID, "data", update.CallbackData())

	user, err := b.checkAuth(ctx, clientID)
	if err != nil {
		return err
	}

	ctx = userutils.NewContextWithUser(ctx, user)

	return b.handleInlineMessage(ctx, update)
}

func (b *Bot) handleInlineMessage(ctx context.Context, update *tgbotapi.Update) error {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.CallbackConfig, error) {
		gameID, err := strconv.ParseUint(update.CallbackQuery.Data, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("parse uint error: %w", err)
		}

		user := userutils.GetUserFromContext(ctx)
		gamePlayer := model.GamePlayer{
			GameID:       int32(gameID),
			UserID:       maybe.Just(user.ID),
			RegisteredBy: user.ID,
			Degree:       model.Degree(gameplayerpb.Degree_DEGREE_LIKELY),
		}

		err = b.gamePlayersFacade.RegisterPlayer(ctx, gamePlayer)
		if err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
				return &cb, nil
			} else if errors.Is(err, gameplayers.ErrNoFreeSlot) {
				cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(noFreeSlotLexeme)(ctx))
				return &cb, nil
			} else if errors.Is(err, gameplayers.ErrGamePlayerAlreadyRegistered) {
				cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(youAreAlreadyRegisteredForTheGameLexeme)(ctx))
				return &cb, nil
			} else if errors.Is(err, games.ErrGameHasPassed) {
				cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(gameHasPassedLexeme)(ctx))
				return &cb, nil
			} else if errors.Is(err, gameplayers.ErrThereAreNoRegistrationForTheGame) {
				cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(thereAreNoRegistrationForTheGameLexeme)(ctx))
				return &cb, nil
			}

			return nil, err
		}

		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(youAreRegisteredForTheGameLexeme)(ctx))
		return &cb, nil
	}

	cb, err := fn(ctx, update)
	if err != nil {
		return fmt.Errorf("prepare inline message response error: %w", err)
	}

	if _, err := b.bot.Request(cb); err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
	}

	return nil
}
