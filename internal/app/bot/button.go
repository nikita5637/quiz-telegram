package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	"github.com/nikita5637/quiz-telegram/internal/pkg/commands"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/icons"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	callbackdatautils "github.com/nikita5637/quiz-telegram/internal/pkg/utils/callbackdata"
)

var (
	cashGamePaymentLexeme = i18n.Lexeme{
		Key:      "cash_game_payment",
		FallBack: "We play for money",
	}
	freeGamePaymentLexeme = i18n.Lexeme{
		Key:      "free_game_payment",
		FallBack: "We play for free",
	}
	gamePhotosLexeme = i18n.Lexeme{
		Key:      "game_photos",
		FallBack: "Game photos",
	}
	legionerIsLikelyToComeLexeme = i18n.Lexeme{
		Key:      "legioner_is_likely_to_come",
		FallBack: "Legioner is likely to come",
	}
	legionerIsUnlikelyToComeLexeme = i18n.Lexeme{
		Key:      "legioner_is_unlikely_to_come",
		FallBack: "Legioner is unlikely to come",
	}
	legionerWillNotComeLexeme = i18n.Lexeme{
		Key:      "legioner_will_not_come",
		FallBack: "Legioner will not come",
	}
	mixGamePaymentLexeme = i18n.Lexeme{
		Key:      "mix_game_payment",
		FallBack: "Mixed payment type",
	}
	playerIsLikelyToComeLexeme = i18n.Lexeme{
		Key:      "player_is_likely_to_come",
		FallBack: "I am likely to come",
	}
	playerIsUnlikelyToComeLexeme = i18n.Lexeme{
		Key:      "player_is_unlikely_to_come",
		FallBack: "I am unlikely to come",
	}
	playerWillNotComeLexeme = i18n.Lexeme{
		Key:      "player_will_not_come",
		FallBack: "I will not come",
	}
	registeredGameLexeme = i18n.Lexeme{
		Key:      "registered_game",
		FallBack: "We are registered for the game",
	}
	registerForLotteryLexeme = i18n.Lexeme{
		Key:      "register_for_lottery",
		FallBack: "Register for lottery",
	}
	registrationLink = i18n.Lexeme{
		Key:      "registration_link",
		FallBack: "Registration link",
	}
	routeToBarLexeme = i18n.Lexeme{
		Key:      "route_to_bar",
		FallBack: "Route to bar",
	}
	unregisteredGameLexeme = i18n.Lexeme{
		Key:      "unregistered_game",
		FallBack: "We are unregistered for the game",
	}
)

func (b *Bot) gamePhotosButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.GetGamePhotosData{
		GameID: gameID,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGamePhotos, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	return tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.Photo, i18n.GetTranslator(gamePhotosLexeme)(ctx)),
		CallbackData: &callbackData,
	}, nil
}

func (b *Bot) lotteryButton(ctx context.Context, gameID int32, leagueID int32, rootGamesListCommand commands.Command) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.LotteryData{
		GameID:                  gameID,
		GetRootGamesListCommand: rootGamesListCommand,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandLottery, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.Lottery, i18n.GetTranslator(registerForLotteryLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	if leagueID == model.LeagueSquiz {
		text := fmt.Sprintf("%s %s", icons.Lottery, i18n.GetTranslator(registrationLink)(ctx))
		btn = tgbotapi.NewInlineKeyboardButtonURL(text, "https://spb.squiz.ru/game")
	}

	return btn, nil
}

func (b *Bot) nextPaymentButton(ctx context.Context, gameID int32, currentPayment int32, rootGamesListCommand commands.Command) (tgbotapi.InlineKeyboardButton, error) {
	nextPayment := int32(0)
	text := ""
	switch currentPayment {
	case int32(gamepb.Payment_PAYMENT_CASH):
		nextPayment = int32(gamepb.Payment_PAYMENT_CERTIFICATE)
		text = fmt.Sprintf("%s %s", icons.FreeGamePayment, i18n.GetTranslator(freeGamePaymentLexeme)(ctx))
	case int32(gamepb.Payment_PAYMENT_CERTIFICATE):
		nextPayment = int32(gamepb.Payment_PAYMENT_MIXED)
		text = fmt.Sprintf("%s %s", icons.MixGamePayment, i18n.GetTranslator(mixGamePaymentLexeme)(ctx))
	case int32(gamepb.Payment_PAYMENT_MIXED):
		nextPayment = int32(gamepb.Payment_PAYMENT_CASH)
		text = fmt.Sprintf("%s %s", icons.CashGamePayment, i18n.GetTranslator(cashGamePaymentLexeme)(ctx))
	}

	payload := &commands.UpdatePaymentData{
		GameID:                  gameID,
		Payment:                 int32(nextPayment),
		GetRootGamesListCommand: rootGamesListCommand,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandUpdatePayment, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	btn := tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) playersListButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.PlayersListByGameData{
		GameID: gameID,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandPlayersListByGame, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	return tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.ListOfPlayers, i18n.GetTranslator(listOfPlayersLexeme)(ctx)),
		CallbackData: &callbackData,
	}, nil
}

func (b *Bot) registerGameButton(ctx context.Context, gameID int32, rootGamesListCommand commands.Command) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.RegisterGameData{
		GameID:                  gameID,
		GetRootGamesListCommand: rootGamesListCommand,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandRegisterGame, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	return tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.RegisteredGame, i18n.GetTranslator(registeredGameLexeme)(ctx)),
		CallbackData: &callbackData,
	}, nil
}

func (b *Bot) registerPlayerButton(ctx context.Context, gameID, userID, registeredBy int32, degree model.Degree, rootGamesListCommand commands.Command) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.RegisterPlayerData{
		GameID:                  gameID,
		UserID:                  userID,
		RegisteredBy:            registeredBy,
		Degree:                  degree,
		GetRootGamesListCommand: rootGamesListCommand,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandRegisterPlayer, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	text := ""
	if userID == registeredBy && degree == model.DegreeLikely {
		text = fmt.Sprintf("%s %s", icons.PlayerLikely, i18n.GetTranslator(playerIsLikelyToComeLexeme)(ctx))
	}
	if userID == registeredBy && degree == model.DegreeUnlikely {
		text = fmt.Sprintf("%s %s", icons.PlayerUnlikely, i18n.GetTranslator(playerIsUnlikelyToComeLexeme)(ctx))
	}
	if userID != registeredBy && degree == model.DegreeLikely {
		text = fmt.Sprintf("%s %s", icons.LegionerLikely, i18n.GetTranslator(legionerIsLikelyToComeLexeme)(ctx))
	}
	if userID != registeredBy && degree == model.DegreeUnlikely {
		text = fmt.Sprintf("%s %s", icons.LegionerUnlikely, i18n.GetTranslator(legionerIsUnlikelyToComeLexeme)(ctx))
	}
	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}, nil
}

func (b *Bot) unregisterGameButton(ctx context.Context, gameID int32, rootGamesListCommand commands.Command) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.UnregisterGameData{
		GameID:                  gameID,
		GetRootGamesListCommand: rootGamesListCommand,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandUnregisterGame, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	return tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.UnregisteredGame, i18n.GetTranslator(unregisteredGameLexeme)(ctx)),
		CallbackData: &callbackData,
	}, nil
}

func (b *Bot) unregisterPlayerButton(ctx context.Context, gameID, userID, registeredBy int32, rootGamesListCommand commands.Command) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.UnregisterPlayerData{
		GameID:                  gameID,
		UserID:                  userID,
		RegisteredBy:            registeredBy,
		Degree:                  model.DegreeInvalid,
		GetRootGamesListCommand: rootGamesListCommand,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandUnregisterPlayer, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	text := fmt.Sprintf("%s %s", icons.LegionerWillNotCome, i18n.GetTranslator(legionerWillNotComeLexeme)(ctx))
	if userID == registeredBy {
		text = fmt.Sprintf("%s %s", icons.PlayerWillNotCome, i18n.GetTranslator(playerWillNotComeLexeme)(ctx))
	}

	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}, nil
}

func (b *Bot) updatePlayerRegistionButton(ctx context.Context, gameID, userID, registeredBy int32, degree model.Degree, rootGamesListCommand commands.Command) (tgbotapi.InlineKeyboardButton, error) {
	payload := &commands.UpdatePlayerRegistration{
		GameID:                  gameID,
		UserID:                  userID,
		RegisteredBy:            registeredBy,
		Degree:                  degree,
		GetRootGamesListCommand: rootGamesListCommand,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandUpdatePlayerRegistration, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	text := fmt.Sprintf("%s %s", icons.PlayerLikely, i18n.GetTranslator(playerIsLikelyToComeLexeme)(ctx))
	if degree == model.DegreeUnlikely {
		text = fmt.Sprintf("%s %s", icons.PlayerUnlikely, i18n.GetTranslator(playerIsUnlikelyToComeLexeme)(ctx))
	}

	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}, nil
}

func (b *Bot) venueButton(ctx context.Context, placeID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := commands.GetVenueData{
		PlaceID: placeID,
	}

	callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetVenue, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, fmt.Errorf("getting callback data error: %w", err)
	}

	text := fmt.Sprintf("%s %s", icons.Route, i18n.GetTranslator(routeToBarLexeme)(ctx))

	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}, nil
}
