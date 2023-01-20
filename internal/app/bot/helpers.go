package bot

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	uuid_utils "github.com/nikita5637/quiz-telegram/utils/uuid"
)

const (
	cashGamePaymentIcon     = "💵"
	freeGamePaymentIcon     = "🆓"
	legionerLikelyIcon      = "💁‍"
	legionerUnlikelyIcon    = "💁‍"
	legionerWillNotComeIcon = "🙅‍"
	listOfPlayersIcon       = "📑"
	lotteryIcon             = "🍀"
	mixGamePaymentIcon      = "❓"
	playerLikelyIcon        = "🙋‍"
	playerUnlikelyIcon      = "🤷‍"
	playerWillNotComeIcon   = "🙅‍"
	prevPageIcon            = "◀️"
	routeIcon               = "🗺"
)

var (
	backToTheGamesListLexeme = i18n.Lexeme{
		Key:      "back_to_the_games_list",
		FallBack: "Back to the games list",
	}
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
	_, err := b.registratorServiceClient.GetUserByTelegramID(ctx, &registrator.GetUserByTelegramIDRequest{
		TelegramId: clientID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) gamesListButton(ctx context.Context) (tgbotapi.InlineKeyboardButton, error) {
	pbReq := &registrator.GetGamesRequest{
		Active: true,
	}

	request, err := getRequest(ctx, CommandGamesList, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	callbackData := b.registerRequest(ctx, request)

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", prevPageIcon, getTranslator(backToTheGamesListLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) getGameMenu(ctx context.Context, game model.Game) (tgbotapi.InlineKeyboardMarkup, error) {
	var err error

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	if game.NumberLegioners+game.NumberPlayers == game.MaxPlayers {
		if game.My {
			var btn1 tgbotapi.InlineKeyboardButton
			btn1, err = b.unregisterPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_MAIN)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1))
		}

		if game.MyLegioners > 0 {
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
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1), tgbotapi.NewInlineKeyboardRow(btn2))
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

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1), tgbotapi.NewInlineKeyboardRow(btn2))
		if game.MyLegioners > 0 {
			var btn3 tgbotapi.InlineKeyboardButton
			btn3, err = b.unregisterPlayerButton(ctx, game.ID, registrator.PlayerType_PLAYER_TYPE_LEGIONER)
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, err
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn3))
		}
	}

	if game.NumberLegioners+game.NumberPlayers > 0 {
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

		var btnNextPayment tgbotapi.InlineKeyboardButton
		btnNextPayment, err = b.nextPaymentButton(ctx, game.ID, game.Payment)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNextPayment))
	}

	if game.WithLottery {
		var btnLottery tgbotapi.InlineKeyboardButton
		btnLottery, err = b.lotteryButton(ctx, game.ID)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnLottery))
	}

	if game.Place.Latitude != 0 && game.Place.Longitude != 0 {
		var btnVenue tgbotapi.InlineKeyboardButton
		btnVenue, err = b.venueButton(ctx, game.Place.Name, game.Place.Address, game.Place.Latitude, game.Place.Longitude)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnVenue))
	}

	var btnGamesList tgbotapi.InlineKeyboardButton
	btnGamesList, err = b.gamesListButton(ctx)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnGamesList))

	return tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}

func (b *Bot) lotteryButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	pbReq := &registrator.RegisterForLotteryRequest{
		GameId: gameID,
	}

	request, err := getRequest(ctx, CommandLottery, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}
	callbackData := b.registerRequest(ctx, request)

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", lotteryIcon, getTranslator(registerForLotteryLexeme)(ctx)),
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
		text = fmt.Sprintf("%s %s :)", freeGamePaymentIcon, getTranslator(freeGamePaymentLexeme)(ctx))
	case model.PaymentTypeCertificate:
		nextPayment = model.PaymentTypeMixed
		text = fmt.Sprintf("%s %s :|", mixGamePaymentIcon, getTranslator(mixGamePaymentLexeme)(ctx))
	case model.PaymentTypeMixed:
		nextPayment = model.PaymentTypeCash
		text = fmt.Sprintf("%s %s :(", cashGamePaymentIcon, getTranslator(cashGamePaymentLexeme)(ctx))
	}

	pbReq := &registrator.UpdatePaymentRequest{
		GameId:  gameID,
		Payment: registrator.Payment(nextPayment),
	}

	request, err := getRequest(ctx, CommandUpdatePayment, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	callbackData := b.registerRequest(ctx, request)

	btn := tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) playersListButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	pbReq := &registrator.GetPlayersByGameIDRequest{
		GameId: gameID,
	}

	request, err := getRequest(ctx, CommandPlayersListByGame, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	callbackData := b.registerRequest(ctx, request)

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s", listOfPlayersIcon, getTranslator(listOfPlayersLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) registerGameButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	pbReq := &registrator.RegisterGameRequest{
		GameId: gameID,
	}

	request, err := getRequest(ctx, CommandRegisterGame, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	callbackData := b.registerRequest(ctx, request)

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s :)", registeredGameIcon, getTranslator(registeredGameLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) registerPlayerButton(ctx context.Context, gameID int32, playerType registrator.PlayerType, degree registrator.Degree) (tgbotapi.InlineKeyboardButton, error) {
	pbReq := &registrator.RegisterPlayerRequest{
		GameId:     gameID,
		PlayerType: playerType,
		Degree:     degree,
	}

	request, err := getRequest(ctx, CommandRegisterPlayer, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	callbackData := b.registerRequest(ctx, request)

	text := ""
	if playerType == registrator.PlayerType_PLAYER_TYPE_MAIN && degree == registrator.Degree_DEGREE_LIKELY {
		text = fmt.Sprintf("%s %s :)", playerLikelyIcon, getTranslator(playerIsLikelyToComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_MAIN && degree == registrator.Degree_DEGREE_UNLIKELY {
		text = fmt.Sprintf("%s %s :|", playerUnlikelyIcon, getTranslator(playerIsUnlikelyToComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_LEGIONER && degree == registrator.Degree_DEGREE_LIKELY {
		text = fmt.Sprintf("%s %s :)", legionerLikelyIcon, getTranslator(legionerIsLikelyToComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_LEGIONER && degree == registrator.Degree_DEGREE_UNLIKELY {
		text = fmt.Sprintf("%s %s :|", legionerUnlikelyIcon, getTranslator(legionerIsUnlikelyToComeLexeme)(ctx))
	}
	btn := tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) registerRequest(ctx context.Context, request model.Request) string {
	var err error
	var requestUUID string
	for i := 0; i < 10; i++ {
		requestUUID, err = b.requestsFacade.RegisterRequest(ctx, request)
		if err != nil {
			continue
		}

		logger.DebugKV(ctx, "registered new request", "uuid", requestUUID, "groupUUID", request.GroupUUID, "body", string(request.Body))
		return requestUUID
	}

	logger.Errorf(ctx, "register request error: %w", err)
	return ""
}

func (b *Bot) unregisterGameButton(ctx context.Context, gameID int32) (tgbotapi.InlineKeyboardButton, error) {
	pbReq := &registrator.UnregisterGameRequest{
		GameId: gameID,
	}

	request, err := getRequest(ctx, CommandUnregisterGame, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	callbackData := b.registerRequest(ctx, request)

	btn := tgbotapi.InlineKeyboardButton{
		Text:         fmt.Sprintf("%s %s :(", unregisteredGameIcon, getTranslator(unregisteredGameLexeme)(ctx)),
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) unregisterPlayerButton(ctx context.Context, gameID int32, playerType registrator.PlayerType) (tgbotapi.InlineKeyboardButton, error) {
	pbReq := &registrator.UnregisterPlayerRequest{
		GameId:     gameID,
		PlayerType: playerType,
	}

	request, err := getRequest(ctx, CommandUnregisterPlayer, pbReq)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}
	callbackData := b.registerRequest(ctx, request)

	text := ""
	if playerType == registrator.PlayerType_PLAYER_TYPE_MAIN {
		text = fmt.Sprintf("%s %s :(", playerWillNotComeIcon, getTranslator(playerWillNotComeLexeme)(ctx))
	}
	if playerType == registrator.PlayerType_PLAYER_TYPE_LEGIONER {
		text = fmt.Sprintf("%s %s :(", legionerWillNotComeIcon, getTranslator(legionerWillNotComeLexeme)(ctx))
	}
	btn := tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}

	return btn, nil
}

func (b *Bot) unregisterRequest(ctx context.Context, uuid string) error {
	err := b.requestsFacade.UnregisterRequest(ctx, uuid)
	if err != nil {
		logger.Errorf(ctx, "unregister request error: %w", err)
		return err
	}

	logger.DebugKV(ctx, "unregistered request", "uuid", uuid)
	return nil
}

func (b *Bot) venueButton(ctx context.Context, title, address string, latitude, longitude float32) (tgbotapi.InlineKeyboardButton, error) {
	venue := VenueData{
		Title:     title,
		Address:   address,
		Latitude:  latitude,
		Longitude: longitude,
	}

	request, err := getRequest(ctx, CommandGetVenue, venue)
	if err != nil {
		return tgbotapi.InlineKeyboardButton{}, err
	}

	text := fmt.Sprintf("%s %s", routeIcon, getTranslator(routeToBarLexeme)(ctx))
	callbackData := b.registerRequest(ctx, request)

	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &callbackData,
	}, nil
}

func convertPBGameToModelGame(pbGame *registrator.Game) model.Game {
	ret := model.Game{
		Date:       model.DateTime(pbGame.GetDate().AsTime()),
		GameType:   int32(pbGame.GetType()),
		ID:         pbGame.GetId(),
		ExternalID: pbGame.GetExternalId(),
		LeagueID:   pbGame.GetLeagueId(),
		MaxPlayers: byte(pbGame.GetMaxPlayers()),
		Number:     pbGame.GetNumber(),
		PlaceID:    pbGame.GetPlaceId(),
		Registered: pbGame.GetRegistered(),
		Payment:    model.PaymentType(pbGame.GetPayment()),
	}

	ret.My = pbGame.GetMy()
	ret.MyLegioners = byte(pbGame.GetNumberOfMyLegioners())
	ret.NumberLegioners = byte(pbGame.GetNumberOfLegioners())
	ret.NumberPlayers = byte(pbGame.GetNumberOfPlayers())

	return ret
}

func getTranslator(lexeme i18n.Lexeme) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return i18n.Translate(ctx, lexeme.Key, lexeme.FallBack)
	}
}

func getRequest[T any](ctx context.Context, command Command, pbReq T) (model.Request, error) {
	body, err := json.Marshal(pbReq)
	if err != nil {
		return model.Request{}, err
	}

	req := TelegramRequest{
		Command: command,
		Body:    body,
	}

	requestBody, err := json.Marshal(req)
	if err != nil {
		return model.Request{}, err
	}

	request := model.Request{
		GroupUUID: uuid_utils.GroupUUIDFromContext(ctx),
		Body:      requestBody,
	}

	return request, nil
}

func convertPBPlaceToModelPlace(place *registrator.Place) model.Place {
	return model.Place{
		ID:        place.GetId(),
		Address:   place.GetAddress(),
		Name:      place.GetName(),
		ShortName: place.GetShortName(),
		Longitude: place.GetLongitude(),
		Latitude:  place.GetLatitude(),
		MenuLink:  place.GetMenuLink(),
	}
}
