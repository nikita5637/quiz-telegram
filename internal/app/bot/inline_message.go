package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mono83/maybe"
	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gameplayers"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	userutils "github.com/nikita5637/quiz-telegram/utils/user"
)

func (b *Bot) handleInlineMessage(ctx context.Context, update *tgbotapi.Update) error {
	user := userutils.GetUserFromContext(ctx)
	logger.DebugKV(ctx, "new inline message incoming", "user", user, "data", update.CallbackData())

	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.CallbackConfig, error) {
		callbackQueryID := update.CallbackQuery.ID

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

		cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreSignedUpForTheGameLexeme)(ctx))
		if err = b.gamePlayersFacade.RegisterPlayer(ctx, gamePlayer); err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
			} else if errors.Is(err, gameplayers.ErrNoFreeSlot) {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(noFreeSlotLexeme)(ctx))
			} else if errors.Is(err, gameplayers.ErrGamePlayerAlreadyRegistered) {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreAlreadyRegisteredForTheGameLexeme)(ctx))
			} else if errors.Is(err, games.ErrGameHasPassed) {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(games.GameHasPassedLexeme)(ctx))
			} else if errors.Is(err, gameplayers.ErrThereAreNoRegistrationForTheGame) {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(thereAreNoRegistrationForTheGameLexeme)(ctx))
			} else {
				return nil, err
			}
		}

		return &cb, nil
	}

	cb, err := fn(ctx, update)
	if err != nil {
		return fmt.Errorf("preparing inline message callback error: %w", err)
	}

	if _, err = b.bot.Request(cb); err != nil {
		return fmt.Errorf("sending callback error: %w", err)
	}

	return nil
}
