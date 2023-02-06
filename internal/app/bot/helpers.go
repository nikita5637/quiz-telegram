package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/icons"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
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
	listOfPlayersLexeme = i18n.Lexeme{
		Key:      "list_of_players",
		FallBack: "List of players",
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
	registerForLotteryLexeme = i18n.Lexeme{
		Key:      "register_for_lottery",
		FallBack: "Register for lottery",
	}
	routeToBarLexeme = i18n.Lexeme{
		Key:      "route_to_bar",
		FallBack: "Route to bar",
	}
)

func (b *Bot) checkAuth(ctx context.Context, clientID int64) error {
	_, err := b.usersFacade.GetUserByTelegramID(ctx, clientID)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) getGameMenu(ctx context.Context, game model.Game, page uint32) (tgbotapi.InlineKeyboardMarkup, error) {
	switch page {
	case 0:
		return b.getGameMenuFirstPage(ctx, game)
	case 1:
		return b.getGameMenuSecondPage(ctx, game)
	}

	return tgbotapi.InlineKeyboardMarkup{}, nil
}

func (b *Bot) getGameMenuFirstPage(ctx context.Context, game model.Game) (tgbotapi.InlineKeyboardMarkup, error) {
	var err error

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	if game.WithLottery {
		var btnLottery tgbotapi.InlineKeyboardButton
		btnLottery, err = b.lotteryButton(ctx, game.ID)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnLottery))
	}

	if game.NumberOfLegioners+game.NumberOfPlayers == game.MaxPlayers {
		if game.My {
			var btn1 tgbotapi.InlineKeyboardButton
			btn1, err = b.unregisterPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_MAIN)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1))
		}

		if game.NumberOfMyLegioners > 0 {
			var btn1 tgbotapi.InlineKeyboardButton
			btn1, err = b.unregisterPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_LEGIONER)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1))
		}
	} else {
		if game.My {
			var btn1 tgbotapi.InlineKeyboardButton
			btn1, err = b.unregisterPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_MAIN)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1))
		} else {
			var btn1 tgbotapi.InlineKeyboardButton
			btn1, err = b.registerPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_MAIN, registrator.Degree_DEGREE_LIKELY)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			var btn2 tgbotapi.InlineKeyboardButton
			btn2, err = b.registerPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_MAIN, registrator.Degree_DEGREE_UNLIKELY)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1, btn2))
		}

		var btn1 tgbotapi.InlineKeyboardButton
		btn1, err = b.registerPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_LEGIONER, registrator.Degree_DEGREE_LIKELY)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}

		var btn2 tgbotapi.InlineKeyboardButton
		btn2, err = b.registerPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_LEGIONER, registrator.Degree_DEGREE_UNLIKELY)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1, btn2))
		if game.NumberOfMyLegioners > 0 {
			var btn3 tgbotapi.InlineKeyboardButton
			btn3, err = b.unregisterPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_LEGIONER)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn3))
		}
	}

	if game.NumberOfLegioners+game.NumberOfPlayers > 0 {
		var btnPlayersList tgbotapi.InlineKeyboardButton
		btnPlayersList, err = b.playersListButton(ctx, game.ID)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPlayersList))
	}

	if !game.Registered {
		var btnRegisterGame tgbotapi.InlineKeyboardButton
		btnRegisterGame, err = b.registerGameButton(ctx, game.ID)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnRegisterGame))
	} else {
		var btnUnregisterGame tgbotapi.InlineKeyboardButton
		btnUnregisterGame, err = b.unregisterGameButton(ctx, game.ID)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnUnregisterGame))
	}

	getGameData := &GetGameData{
		GameID:    game.ID,
		PageIndex: 1,
	}

	var callbackData string
	callbackData, err = getCallbackData(ctx, CommandGetGame, getGameData)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	btnNextMenuPage := tgbotapi.InlineKeyboardButton{
		Text:         icons.NextPage,
		CallbackData: &callbackData,
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNextMenuPage))

	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}

func (b *Bot) getGameMenuSecondPage(ctx context.Context, game model.Game) (tgbotapi.InlineKeyboardMarkup, error) {
	var err error

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	if game.Registered {
		var btnNextPayment tgbotapi.InlineKeyboardButton
		btnNextPayment, err = b.nextPaymentButton(ctx, game.ID, game.Payment)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNextPayment))
	}

	barButtonsRow := []tgbotapi.InlineKeyboardButton{}
	if game.Place.Latitude != 0 && game.Place.Longitude != 0 {
		var btnVenue tgbotapi.InlineKeyboardButton
		btnVenue, err = b.venueButton(ctx, game.Place.ID)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		barButtonsRow = append(barButtonsRow, btnVenue)
	}

	if game.Place.MenuLink != "" {
		btnMenu := tgbotapi.NewInlineKeyboardButtonURL("🍴 Меню ресторана", game.Place.MenuLink)
		barButtonsRow = append(barButtonsRow, btnMenu)
	}

	if len(barButtonsRow) > 0 {
		rows = append(rows, barButtonsRow)
	}

	getGameData := &GetGameData{
		GameID:    game.ID,
		PageIndex: 0,
	}

	var callbackData string
	callbackData, err = getCallbackData(ctx, CommandGetGame, getGameData)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	btnPrevMenuPage := tgbotapi.InlineKeyboardButton{
		Text:         icons.PrevPage,
		CallbackData: &callbackData,
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPrevMenuPage))

	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}

func (b *Bot) lotteryButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := &LotteryData{
		GameID: gameID,
	}

	callbackData, err := getCallbackData(ctx, CommandLottery, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.Lottery, getTranslator(registerForLotteryLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) nextPaymentButton(ctx context.Context, gameID int32, currentPayment model.PaymentType) (tgbotapi.InlineKeyboardButton, error) {
	nextPayment := model.PaymentTypeInvalid
	text := ""
	switch currentPayment {
	case model.PaymentTypeCash:
		nextPayment = model.PaymentTypeCertificate
		text = fmt.Sprintf("%s %s", icons.FreeGamePayment, getTranslator(freeGamePaymentLexeme)(ctx))
	case model.PaymentTypeCertificate:
		nextPayment = model.PaymentTypeMixed
		text = fmt.Sprintf("%s %s", icons.MixGamePayment, getTranslator(mixGamePaymentLexeme)(ctx))
	case model.PaymentTypeMixed:
		nextPayment = model.PaymentTypeCash
		text = fmt.Sprintf("%s %s", icons.CashGamePayment, getTranslator(cashGamePaymentLexeme)(ctx))
	}

	payload := &UpdatePaymentData{
		GameID:  gameID,
		Payment: int32(nextPayment),
	}

	callbackData, err := getCallbackData(ctx, CommandUpdatePayment, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	btn := tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) playersListButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := &PlayersListByGameData{
		GameID: gameID,
	}

	callbackData, err := getCallbackData(ctx, CommandPlayersListByGame, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.ListOfPlayers, getTranslator(listOfPlayersLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) registerGameButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := &RegisterGameData{
		GameID: gameID,
	}

	callbackData, err := getCallbackData(ctx, CommandRegisterGame, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.RegisteredGame, getTranslator(registeredGameLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) registerPlayerButton(ctx context.Context, gameID int32, playerType registrator.PlayerType, degree registrator.Degree) (tgbotapi.InlineKeyboardButton, error) {
	payload := &RegisterPlayerData{
		GameID:     gameID,
		PlayerType: int32(playerType),
		Degree:     int32(degree),
	}

	callbackData, err := getCallbackData(ctx, CommandRegisterPlayer, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	text := ""
	if playerType == registrator.PlayerType_PLAYER_TYPE_MAIN && degree == registrator.Degree_DEGREE_LIKELY {
		text = fmt.Sprintf("%s %s", icons.PlayerLikely, getTranslator(playerIsLikelyToComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_MAIN && degree == registrator.Degree_DEGREE_UNLIKELY {
		text = fmt.Sprintf("%s %s", icons.PlayerUnlikely, getTranslator(playerIsUnlikelyToComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_LEGIONER && degree == registrator.Degree_DEGREE_LIKELY {
		text = fmt.Sprintf("%s %s", icons.LegionerLikely, getTranslator(legionerIsLikelyToComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_LEGIONER && degree == registrator.Degree_DEGREE_UNLIKELY {
		text = fmt.Sprintf("%s %s", icons.LegionerUnlikely, getTranslator(legionerIsUnlikelyToComeLexeme)(ctx))
	}
	btn := tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) unregisterGameButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := &UnregisterGameData{
		GameID: gameID,
	}

	callbackData, err := getCallbackData(ctx, CommandUnregisterGame, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", icons.UnregisteredGame, getTranslator(unregisteredGameLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) unregisterPlayerButton(ctx context.Context, gameID int32, playerType registrator.PlayerType) (tgbotapi.InlineKeyboardButton, error) {
	payload := &UnregisterPlayerData{
		GameID:     gameID,
		PlayerType: int32(playerType),
	}

	callbackData, err := getCallbackData(ctx, CommandUnregisterPlayer, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	text := ""
	if playerType == registrator.PlayerType_PLAYER_TYPE_MAIN {
		text = fmt.Sprintf("%s %s", icons.PlayerWillNotCome, getTranslator(playerWillNotComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_LEGIONER {
		text = fmt.Sprintf("%s %s", icons.LegionerWillNotCome, getTranslator(legionerWillNotComeLexeme)(ctx))
	}
	btn := tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) venueButton(ctx context.Context, placeID int32) (tgbotapi.InlineKeyboardButton, error) {
	payload := GetVenueData{
		PlaceID: placeID,
	}

	callbackData, err := getCallbackData(ctx, CommandGetVenue, payload)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	text := fmt.Sprintf("%s %s", icons.Route, getTranslator(routeToBarLexeme)(ctx))

	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}, nil
}

func getTranslator(lexeme i18n.Lexeme) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return i18n.Translate(ctx, lexeme.Key, lexeme.FallBack)
	}
}

func getCallbackData(ctx context.Context, command Command, payload interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req := TelegramRequest{
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

func replyKeyboardMarkup(ctx context.Context) tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(getTranslator(myGamesLexeme)(ctx)),
			tgbotapi.NewKeyboardButton(getTranslator(registeredGamesLexeme)(ctx)),
		},
		[]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton(getTranslator(settingsLexeme)(ctx)),
		},
	)
	kb.ResizeKeyboard = true

	return kb
}
