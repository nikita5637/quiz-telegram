package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	prevPageStringText   = "<"
	registeredGameIcon   = "✅"
	unregisteredGameIcon = "❌"
)

var (
	degreeMap = map[registrator.Degree]i18n.Lexeme{
		registrator.Degree_DEGREE_INVALID: i18n.Lexeme{
			Key:      "invalid_degree",
			FallBack: "invalid degree",
		},
		registrator.Degree_DEGREE_LIKELY: i18n.Lexeme{
			Key:      "plays_likely",
			FallBack: "plays likely",
		},
		registrator.Degree_DEGREE_UNLIKELY: i18n.Lexeme{
			Key:      "plays_unlikely",
			FallBack: "plays unlikely",
		},
	}

	enterYourEmailLexeme = i18n.Lexeme{
		Key:      "enter_your_email",
		FallBack: "OK. Enter your email.",
	}
	enterYourNameLexeme = i18n.Lexeme{
		Key:      "enter_your_name",
		FallBack: "OK. Enter your name.",
	}
	enterYourPhoneLexeme = i18n.Lexeme{
		Key:      "enter_your_phone",
		FallBack: "OK. Enter your phone(format +79XXXXXXXXX).",
	}
	gameNotFoundLexeme = i18n.Lexeme{
		Key:      "game_not_found",
		FallBack: "Game not found",
	}
	legionerByLexeme = i18n.Lexeme{
		Key:      "legioner_by",
		FallBack: "Legioner by",
	}
	listOfPlayersIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_players_is_empty",
		FallBack: "There are not players",
	}
	registeredGameLexeme = i18n.Lexeme{
		Key:      "registered_game",
		FallBack: "We are registered for the game",
	}
	unregisteredGameLexeme = i18n.Lexeme{
		Key:      "unregistered_game",
		FallBack: "We are unregistered for the game",
	}
	youHaveSuccessfullyRegisteredInLotteryLexeme = i18n.Lexeme{
		Key:      "you_have_successfully_registered_in_the_lottery",
		FallBack: "You have successfully registered in the lottery",
	}
	yourLotteryNumberIsLexeme = i18n.Lexeme{
		Key:      "your_lottery_number_is",
		FallBack: "Your lottery number is",
	}
)

// HandleCallbackQuery ...
func (b *Bot) HandleCallbackQuery(ctx context.Context, update *tgbotapi.Update) error {
	clientID := update.CallbackQuery.From.ID
	ctx = telegram_utils.NewContextWithClientID(ctx, clientID)

	err := b.checkAuth(ctx, clientID)
	if err != nil {
		name := update.CallbackQuery.Message.Chat.FirstName
		if name == "" {
			name = update.CallbackQuery.Message.Chat.UserName
		}

		_, err = b.registratorServiceClient.CreateUser(ctx, &registrator.CreateUserRequest{
			Name:       name,
			TelegramId: clientID,
			State:      registrator.UserState_USER_STATE_WELCOME,
		})
		if err != nil {
			st := status.Convert(err)

			if st.Code() == codes.AlreadyExists {
				return nil
			}

			return err
		}

		welcomeMessage := welcomeMessage(ctx, clientID, name)
		_, err = b.bot.Send(welcomeMessage)
		return err
	}

	telegramRequest := TelegramRequest{}

	err = json.Unmarshal([]byte(update.CallbackData()), &telegramRequest)
	if err != nil {
		return fmt.Errorf("telegram request unmarshaling error: %w", err)
	}

	type handlerFunc func(ctx context.Context) error
	var handler handlerFunc
	switch telegramRequest.Command {
	case CommandChangeEmail:
		handler = func(ctx context.Context) error {
			return b.handleChangeEmail(ctx, update, telegramRequest)
		}
	case CommandChangeName:
		handler = func(ctx context.Context) error {
			return b.handleChangeName(ctx, update, telegramRequest)
		}
	case CommandChangePhone:
		handler = func(ctx context.Context) error {
			return b.handleChangePhone(ctx, update, telegramRequest)
		}
	case CommandGetGamesList:
		handler = func(ctx context.Context) error {
			return b.handleGetGamesList(ctx, update, telegramRequest)
		}
	case CommandGetGame:
		handler = func(ctx context.Context) error {
			return b.handleGetGame(ctx, update, telegramRequest)
		}
	case CommandGetGamePhotos:
		handler = func(ctx context.Context) error {
			return b.handleGetGamePhotos(ctx, update, telegramRequest)
		}
	case CommandGetListGamesWithPhotosNextPage:
		handler = func(ctx context.Context) error {
			return b.handleGetListGamesWithPhotosNextPage(ctx, update, telegramRequest)
		}
	case CommandGetListGamesWithPhotosPrevPage:
		handler = func(ctx context.Context) error {
			return b.handleGetListGamesWithPhotosPrevPage(ctx, update, telegramRequest)
		}
	case CommandGetVenue:
		handler = func(ctx context.Context) error {
			return b.handleGetVenue(ctx, update, telegramRequest)
		}
	case CommandLottery:
		handler = func(ctx context.Context) error {
			return b.handleLottery(ctx, update, telegramRequest)
		}
	case CommandPlayersListByGame:
		handler = func(ctx context.Context) error {
			return b.handlePlayersList(ctx, update, telegramRequest)
		}
	case CommandRegisterGame:
		handler = func(ctx context.Context) error {
			return b.handleRegisterGame(ctx, update, telegramRequest)
		}
	case CommandRegisterPlayer:
		handler = func(ctx context.Context) error {
			return b.handleRegisterPlayer(ctx, update, telegramRequest)
		}
	case CommandUnregisterGame:
		handler = func(ctx context.Context) error {
			return b.handleUnregisterGame(ctx, update, telegramRequest)
		}
	case CommandUnregisterPlayer:
		handler = func(ctx context.Context) error {
			return b.handleUnregisterPlayer(ctx, update, telegramRequest)
		}
	case CommandUpdatePayment:
		handler = func(ctx context.Context) error {
			return b.handleUpdatePayment(ctx, update, telegramRequest)
		}
	}

	err = handler(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleChangeEmail(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	return b.updateUserState(ctx, update, model.UserStateChangingEmail)
}

func (b *Bot) handleChangeName(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	return b.updateUserState(ctx, update, model.UserStateChangingName)
}

func (b *Bot) handleChangePhone(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	return b.updateUserState(ctx, update, model.UserStateChangingPhone)
}

func (b *Bot) handleGetGamesList(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID

	games, err := b.gamesFacade.GetGames(ctx, true)
	if err != nil {
		return err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &GetGameData{
			GameID: game.ID,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, CommandGetGame, payload)
		if err != nil {
			return err
		}

		text := fmt.Sprintf(gameInfoFormatString, game.League.ShortName, game.Number, game.Place.ShortName, game.DateTime())

		if game.My {
			text = myGamePrefix + text
		} else {
			if game.NumberOfLegioners+game.NumberOfPlayers > 0 {
				text = gameWithPlayersPrefix + text
			}
		}

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg := tgbotapi.NewEditMessageText(clientID, update.CallbackQuery.Message.MessageID, getTranslator(chooseGameLexeme)(ctx))
	inlineKeyboarMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ReplyMarkup = &inlineKeyboarMarkup

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "error while sending message: %w", err)
		return err
	}

	return nil
}

func (b *Bot) handleGetGame(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &GetGameData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		}
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &registrator.GetLotteryStatusRequest{
		GameId: game.ID,
	})
	if err != nil {
		logger.Warnf(ctx, "getting lottery status error: %w", err)
	} else {
		game.WithLottery = lotteryResp.GetActive()
	}

	var menu tgbotapi.InlineKeyboardMarkup
	menu, err = b.getGameMenu(ctx, game)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
	msg.ReplyMarkup = &menu

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleGetGamePhotos(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &GetGamePhotosData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	urls, err := b.gamePhotosFacade.GetPhotosByGameID(ctx, data.GameID)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		}

		return err
	}

	deleteConfig := tgbotapi.NewDeleteMessage(clientID, messageID)
	_, err = b.bot.Request(deleteConfig)
	if err != nil {
		return err
	}

	for _, url := range urls {
		msg := tgbotapi.NewMessage(clientID, url)
		_, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleGetListGamesWithPhotosNextPage(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	gamesWithPhotosListLimit := uint32(config.GetValue("GamesWithPhotosListLimit").Uint64())

	data := &GetGamesWithPhotosData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	games, total, err := b.gamePhotosFacade.GetGamesWithPhotos(ctx, data.Limit, data.Offset)
	if err != nil {
		return err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &GetGamePhotosData{
			GameID: game.ID,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, CommandGetGamePhotos, payload)
		if err != nil {
			return err
		}

		text := fmt.Sprintf(gamePhotosInfoFormatString, game.ResultPlace.String(), game.League.ShortName, game.Number, game.Place.ShortName, game.DateTime())

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	navigateButtonsRow := make([]tgbotapi.InlineKeyboardButton, 0, 2)

	if data.Offset > 0 {
		offset := uint32(0)
		if data.Offset > gamesWithPhotosListLimit {
			offset = data.Offset - gamesWithPhotosListLimit
		}

		payload := &GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: offset,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, CommandGetListGamesWithPhotosPrevPage, payload)
		if err != nil {
			return err
		}

		btnPrev := tgbotapi.InlineKeyboardButton{
			Text:         prevPageStringText,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnPrev)
	}

	leftNext := uint32(0)
	if total > (data.Offset + data.Offset) {
		leftNext = total - (data.Offset + data.Limit)
	}

	if leftNext > 0 {
		payload := &GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: data.Offset + data.Limit,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, CommandGetListGamesWithPhotosNextPage, payload)
		if err != nil {
			return err
		}

		btnNext := tgbotapi.InlineKeyboardButton{
			Text:         nextPageStringText,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnNext)
	}

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, navigateButtonsRow)

	msg := tgbotapi.NewEditMessageReplyMarkup(clientID, messageID, inlineKeyboardMarkup)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleGetListGamesWithPhotosPrevPage(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	gamesWithPhotosListLimit := uint32(config.GetValue("GamesWithPhotosListLimit").Uint64())

	data := &GetGamesWithPhotosData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	games, total, err := b.gamePhotosFacade.GetGamesWithPhotos(ctx, data.Limit, data.Offset)
	if err != nil {
		return err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &GetGamePhotosData{
			GameID: game.ID,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, CommandGetGamePhotos, payload)
		if err != nil {
			return err
		}

		text := fmt.Sprintf(gamePhotosInfoFormatString, game.ResultPlace.String(), game.League.ShortName, game.Number, game.Place.ShortName, game.DateTime())

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	navigateButtonsRow := make([]tgbotapi.InlineKeyboardButton, 0, 2)

	if data.Offset > 0 {
		offset := uint32(0)
		if data.Offset > gamesWithPhotosListLimit {
			offset = data.Offset - gamesWithPhotosListLimit
		}

		payload := &GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: offset,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, CommandGetListGamesWithPhotosPrevPage, payload)
		if err != nil {
			return err
		}

		btnPrev := tgbotapi.InlineKeyboardButton{
			Text:         prevPageStringText,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnPrev)
	}

	leftNext := uint32(0)
	if total > (data.Offset + data.Limit) {
		leftNext = total - (data.Offset + data.Limit)
	}

	if leftNext > 0 {
		payload := &GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: data.Offset + data.Limit,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, CommandGetListGamesWithPhotosNextPage, payload)
		if err != nil {
			return err
		}

		btnNext := tgbotapi.InlineKeyboardButton{
			Text:         nextPageStringText,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnNext)
	}

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, navigateButtonsRow)

	msg := tgbotapi.NewEditMessageReplyMarkup(clientID, messageID, inlineKeyboardMarkup)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleGetVenue(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &GetVenueData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: data.PlaceID,
	})
	if err != nil {
		return err
	}

	deleteMessageConfig := tgbotapi.NewDeleteMessage(clientID, messageID)
	_, err = b.bot.Request(deleteMessageConfig)
	if err != nil {
		return err
	}

	venueConfig := tgbotapi.NewVenue(clientID, placeResp.GetPlace().GetName(), placeResp.GetPlace().GetAddress(), float64(placeResp.GetPlace().GetLatitude()), float64(placeResp.GetPlace().GetLongitude()))
	_, err = b.bot.Request(venueConfig)
	return err
}

func (b *Bot) handleLottery(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &LotteryData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	resp, err := b.croupierServiceClient.RegisterForLottery(ctx, &registrator.RegisterForLotteryRequest{
		GameId: data.GameID,
	})
	if err != nil {
		st := status.Convert(err)

		msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(somethingWentWrongLexeme)(ctx))
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.LocalizedMessage:
				localizedMessage := t.GetMessage()
				msg = tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
			}
		}

		_, err = b.bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(youHaveSuccessfullyRegisteredInLotteryLexeme)(ctx))
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	if resp.GetNumber() > 0 {
		msg := tgbotapi.NewMessage(clientID, fmt.Sprintf("%s: %d", getTranslator(yourLotteryNumberIsLexeme)(ctx), resp.GetNumber()))
		newMsg, err := b.bot.Send(msg)
		if err != nil {
			return err
		}

		unpinMessage := tgbotapi.UnpinAllChatMessagesConfig{
			ChatID: clientID,
		}
		_, err = b.bot.Request(unpinMessage)
		if err != nil {
			return err
		}

		pinMessage := tgbotapi.PinChatMessageConfig{
			ChatID:    clientID,
			MessageID: newMsg.MessageID,
		}
		_, err = b.bot.Request(pinMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handlePlayersList(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &PlayersListByGameData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	resp, err := b.registratorServiceClient.GetPlayersByGameID(ctx, &registrator.GetPlayersByGameIDRequest{
		GameId: data.GameID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		}
		return err
	}

	textBuilder := strings.Builder{}
	for i, player := range resp.GetPlayers() {
		playerName := ""
		if player.UserId > 0 {
			var playerResp *registrator.GetUserByIDResponse
			if playerResp, err = b.registratorServiceClient.GetUserByID(ctx, &registrator.GetUserByIDRequest{
				Id: player.UserId,
			}); err != nil {
				return err
			}
			playerName = playerResp.GetUser().GetName()
		} else {
			var playerResp *registrator.GetUserByIDResponse
			if playerResp, err = b.registratorServiceClient.GetUserByID(ctx, &registrator.GetUserByIDRequest{
				Id: player.RegisteredBy,
			}); err != nil {
				return err
			}
			playerName = fmt.Sprintf("%s %s", getTranslator(legionerByLexeme)(ctx), playerResp.GetUser().GetName())
		}

		if player.GetDegree() == registrator.Degree_DEGREE_UNLIKELY {
			textBuilder.WriteString(fmt.Sprintf("%d. %s (%s)\n", i+1, playerName, getTranslator(degreeMap[registrator.Degree_DEGREE_UNLIKELY])(ctx)))
		} else {
			textBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, playerName))
		}
	}

	text := textBuilder.String()
	if text == "" {
		text = fmt.Sprintf("%s :(", getTranslator(listOfPlayersIsEmptyLexeme)(ctx))
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, text)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleRegisterGame(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &RegisterGameData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.RegisterGame(ctx, &registrator.RegisterGameRequest{
		GameId: data.GameID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		}
		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &registrator.GetLotteryStatusRequest{
		GameId: game.ID,
	})
	if err != nil {
		logger.Warnf(ctx, "getting lottery status error: %w", err)
	} else {
		game.WithLottery = lotteryResp.GetActive()
	}

	var menu tgbotapi.InlineKeyboardMarkup
	menu, err = b.getGameMenu(ctx, game)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
	msg.ReplyMarkup = &menu

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleRegisterPlayer(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &RegisterPlayerData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.RegisterPlayer(ctx, &registrator.RegisterPlayerRequest{
		GameId:     data.GameID,
		PlayerType: registrator.PlayerType(data.PlayerType),
		Degree:     registrator.Degree(data.Degree),
	})
	if err != nil {
		st := status.Convert(err)

		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		} else if st.Code() == codes.AlreadyExists {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		}

		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &registrator.GetLotteryStatusRequest{
		GameId: game.ID,
	})
	if err != nil {
		logger.Warnf(ctx, "getting lottery status error: %w", err)
	} else {
		game.WithLottery = lotteryResp.GetActive()
	}

	var menu tgbotapi.InlineKeyboardMarkup
	menu, err = b.getGameMenu(ctx, game)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
	msg.ReplyMarkup = &menu

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleUnregisterGame(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &UnregisterGameData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.UnregisterGame(ctx, &registrator.UnregisterGameRequest{
		GameId: data.GameID,
	})
	if err != nil {
		st := status.Convert(err)

		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		}

		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &registrator.GetLotteryStatusRequest{
		GameId: game.ID,
	})
	if err != nil {
		logger.Warnf(ctx, "getting lottery status error: %w", err)
	} else {
		game.WithLottery = lotteryResp.GetActive()
	}

	var menu tgbotapi.InlineKeyboardMarkup
	menu, err = b.getGameMenu(ctx, game)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
	msg.ReplyMarkup = &menu

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleUnregisterPlayer(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &UnregisterPlayerData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.UnregisterPlayer(ctx, &registrator.UnregisterPlayerRequest{
		GameId:     data.GameID,
		PlayerType: registrator.PlayerType(data.PlayerType),
	})
	if err != nil {
		st := status.Convert(err)

		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		}

		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &registrator.GetLotteryStatusRequest{
		GameId: game.ID,
	})
	if err != nil {
		logger.Warnf(ctx, "getting lottery status error: %w", err)
	} else {
		game.WithLottery = lotteryResp.GetActive()
	}

	var menu tgbotapi.InlineKeyboardMarkup
	menu, err = b.getGameMenu(ctx, game)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
	msg.ReplyMarkup = &menu

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleUpdatePayment(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &UpdatePaymentData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.UpdatePayment(ctx, &registrator.UpdatePaymentRequest{
		GameId:  data.GameID,
		Payment: registrator.Payment(data.Payment),
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					_, err = b.bot.Send(msg)
					return err
				}
			}
		}
		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &registrator.GetLotteryStatusRequest{
		GameId: game.ID,
	})
	if err != nil {
		logger.Warnf(ctx, "getting lottery status error: %w", err)
	} else {
		game.WithLottery = lotteryResp.GetActive()
	}

	var menu tgbotapi.InlineKeyboardMarkup
	menu, err = b.getGameMenu(ctx, game)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
	msg.ReplyMarkup = &menu

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) updateUserState(ctx context.Context, update *tgbotapi.Update, state model.UserState) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	resp, err := b.registratorServiceClient.GetUserByTelegramID(ctx, &registrator.GetUserByTelegramIDRequest{
		TelegramId: clientID,
	})
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.UpdateUserState(ctx, &registrator.UpdateUserStateRequest{
		UserId: resp.GetUser().GetId(),
		State:  registrator.UserState(state),
	})
	if err != nil {
		return err
	}

	msg := tgbotapi.EditMessageTextConfig{}
	switch state {
	case model.UserStateChangingEmail:
		msg = tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(enterYourEmailLexeme)(ctx))
	case model.UserStateChangingName:
		msg = tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(enterYourNameLexeme)(ctx))
	case model.UserStateChangingPhone:
		msg = tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(enterYourPhoneLexeme)(ctx))
	}

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	return nil
}

func detailInfo(ctx context.Context, game model.Game) string {
	info := strings.Builder{}
	registerStatus := fmt.Sprintf("%s %s :(", unregisteredGameIcon, getTranslator(unregisteredGameLexeme)(ctx))
	if game.Registered {
		registerStatus = fmt.Sprintf("%s %s :)", registeredGameIcon, getTranslator(registeredGameLexeme)(ctx))
	}

	info.WriteString(registerStatus + "\n")

	if game.Payment != model.PaymentTypeInvalid {
		paymentStatus := "❓ Оплата: микс"
		if game.Payment == model.PaymentTypeCash {
			paymentStatus = "💵 Оплата: денюжкой :("
		} else if game.Payment == model.PaymentTypeCertificate {
			paymentStatus = "🆓 Оплата: сертификатом :)"
		}

		info.WriteString(paymentStatus + "\n")
	}

	info.WriteString("#️⃣ Номер пакета: " + game.Number + "\n")
	info.WriteString("📍 Адрес: " + game.Place.Address + "\n")
	info.WriteString("📅 Дата и время: " + game.DateTime().String() + "\n")
	info.WriteString(fmt.Sprintf("👥 Количество игроков: %d/%d/%d", game.NumberOfPlayers, game.NumberOfLegioners, game.MaxPlayers))

	return info.String()
}
