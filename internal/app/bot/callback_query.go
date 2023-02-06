package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/commands"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/icons"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	addressLexeme = i18n.Lexeme{
		Key:      "address",
		FallBack: "Address",
	}
	cardLexeme = i18n.Lexeme{
		Key:      "card",
		FallBack: "Card",
	}
	cashLexeme = i18n.Lexeme{
		Key:      "cash",
		FallBack: "Cash",
	}
	certificateLexeme = i18n.Lexeme{
		Key:      "certificate",
		FallBack: "Certificate",
	}
	dateTimeLexeme = i18n.Lexeme{
		Key:      "datetime",
		FallBack: "Datetime",
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
	gameCostLexeme = i18n.Lexeme{
		Key:      "game_cost",
		FallBack: "Game cost",
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
	mixLexeme = i18n.Lexeme{
		Key:      "mix",
		FallBack: "Mix",
	}
	numberLexeme = i18n.Lexeme{
		Key:      "number",
		FallBack: "Number",
	}
	numberOfPlayersLexeme = i18n.Lexeme{
		Key:      "number_of_players",
		FallBack: "Number of players",
	}
	paymentLexeme = i18n.Lexeme{
		Key:      "payment",
		FallBack: "Payment",
	}
	registeredGameLexeme = i18n.Lexeme{
		Key:      "registered_game",
		FallBack: "We are registered for the game",
	}
	titleLexeme = i18n.Lexeme{
		Key:      "title",
		FallBack: "Title",
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

		_, err = b.usersFacade.CreateUser(ctx, name, clientID, int32(registrator.UserState_USER_STATE_WELCOME))
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

	telegramRequest := commands.TelegramRequest{}

	err = json.Unmarshal([]byte(update.CallbackData()), &telegramRequest)
	if err != nil {
		return fmt.Errorf("telegram request unmarshaling error: %w", err)
	}

	type handlerFunc func(ctx context.Context) error
	var handler handlerFunc
	switch telegramRequest.Command {
	case commands.CommandChangeEmail:
		handler = func(ctx context.Context) error {
			return b.handleChangeEmail(ctx, update, telegramRequest)
		}
	case commands.CommandChangeName:
		handler = func(ctx context.Context) error {
			return b.handleChangeName(ctx, update, telegramRequest)
		}
	case commands.CommandChangePhone:
		handler = func(ctx context.Context) error {
			return b.handleChangePhone(ctx, update, telegramRequest)
		}
	case commands.CommandGetGamesList:
		handler = func(ctx context.Context) error {
			return b.handleGetGamesList(ctx, update, telegramRequest)
		}
	case commands.CommandGetGame:
		handler = func(ctx context.Context) error {
			return b.handleGetGame(ctx, update, telegramRequest)
		}
	case commands.CommandGetGamePhotos:
		handler = func(ctx context.Context) error {
			return b.handleGetGamePhotos(ctx, update, telegramRequest)
		}
	case commands.CommandGetListGamesWithPhotosNextPage:
		handler = func(ctx context.Context) error {
			return b.handleGetListGamesWithPhotosNextPage(ctx, update, telegramRequest)
		}
	case commands.CommandGetListGamesWithPhotosPrevPage:
		handler = func(ctx context.Context) error {
			return b.handleGetListGamesWithPhotosPrevPage(ctx, update, telegramRequest)
		}
	case commands.CommandGetVenue:
		handler = func(ctx context.Context) error {
			return b.handleGetVenue(ctx, update, telegramRequest)
		}
	case commands.CommandLottery:
		handler = func(ctx context.Context) error {
			return b.handleLottery(ctx, update, telegramRequest)
		}
	case commands.CommandPlayersListByGame:
		handler = func(ctx context.Context) error {
			return b.handlePlayersList(ctx, update, telegramRequest)
		}
	case commands.CommandRegisterGame:
		handler = func(ctx context.Context) error {
			return b.handleRegisterGame(ctx, update, telegramRequest)
		}
	case commands.CommandRegisterPlayer:
		handler = func(ctx context.Context) error {
			return b.handleRegisterPlayer(ctx, update, telegramRequest)
		}
	case commands.CommandUnregisterGame:
		handler = func(ctx context.Context) error {
			return b.handleUnregisterGame(ctx, update, telegramRequest)
		}
	case commands.CommandUnregisterPlayer:
		handler = func(ctx context.Context) error {
			return b.handleUnregisterPlayer(ctx, update, telegramRequest)
		}
	case commands.CommandUpdatePayment:
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

func (b *Bot) handleChangeEmail(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(registrator.UserState_USER_STATE_CHANGING_EMAIL))
}

func (b *Bot) handleChangeName(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(registrator.UserState_USER_STATE_CHANGING_NAME))
}

func (b *Bot) handleChangePhone(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(registrator.UserState_USER_STATE_CHANGING_PHONE))
}

func (b *Bot) handleGetGamesList(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID

	games, err := b.gamesFacade.GetGames(ctx, true)
	if err != nil {
		return err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &commands.GetGameData{
			GameID:    game.ID,
			PageIndex: 0,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, commands.CommandGetGame, payload)
		if err != nil {
			return err
		}

		text := fmt.Sprintf(gameInfoFormatString, game.League.ShortName, game.Number, game.Place.ShortName, game.DateTime())

		if game.My {
			text = fmt.Sprintf("%s %s", icons.MyGame, text)
		} else {
			if game.NumberOfLegioners+game.NumberOfPlayers > 0 {
				text = fmt.Sprintf("%s %s", icons.GameWithPlayers, text)
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

func (b *Bot) handleGetGame(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.GetGameData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
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
	menu, err = b.getGameMenu(ctx, game, data.PageIndex)
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

func (b *Bot) handleGetGamePhotos(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.GetGamePhotosData{}

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

func (b *Bot) handleGetListGamesWithPhotosNextPage(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	gamesWithPhotosListLimit := uint32(config.GetValue("GamesWithPhotosListLimit").Uint64())

	data := &commands.GetGamesWithPhotosData{}

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
		payload := &commands.GetGamePhotosData{
			GameID: game.ID,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, commands.CommandGetGamePhotos, payload)
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

		payload := &commands.GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: offset,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, commands.CommandGetListGamesWithPhotosPrevPage, payload)
		if err != nil {
			return err
		}

		btnPrev := tgbotapi.InlineKeyboardButton{
			Text:         icons.PrevPage,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnPrev)
	}

	leftNext := uint32(0)
	if total > (data.Offset + data.Offset) {
		leftNext = total - (data.Offset + data.Limit)
	}

	if leftNext > 0 {
		payload := &commands.GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: data.Offset + data.Limit,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, commands.CommandGetListGamesWithPhotosNextPage, payload)
		if err != nil {
			return err
		}

		btnNext := tgbotapi.InlineKeyboardButton{
			Text:         icons.NextPage,
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

func (b *Bot) handleGetListGamesWithPhotosPrevPage(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	gamesWithPhotosListLimit := uint32(config.GetValue("GamesWithPhotosListLimit").Uint64())

	data := &commands.GetGamesWithPhotosData{}

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
		payload := &commands.GetGamePhotosData{
			GameID: game.ID,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, commands.CommandGetGamePhotos, payload)
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

		payload := &commands.GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: offset,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, commands.CommandGetListGamesWithPhotosPrevPage, payload)
		if err != nil {
			return err
		}

		btnPrev := tgbotapi.InlineKeyboardButton{
			Text:         icons.PrevPage,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnPrev)
	}

	leftNext := uint32(0)
	if total > (data.Offset + data.Limit) {
		leftNext = total - (data.Offset + data.Limit)
	}

	if leftNext > 0 {
		payload := &commands.GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: data.Offset + data.Limit,
		}

		var callbackData string
		callbackData, err = getCallbackData(ctx, commands.CommandGetListGamesWithPhotosNextPage, payload)
		if err != nil {
			return err
		}

		btnNext := tgbotapi.InlineKeyboardButton{
			Text:         icons.NextPage,
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

func (b *Bot) handleGetVenue(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID

	data := &commands.GetVenueData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	place, err := b.placesFacade.GetPlaceByID(ctx, data.PlaceID)
	if err != nil {
		return err
	}

	venueConfig := tgbotapi.NewVenue(clientID, place.Name, place.Address, float64(place.Latitude), float64(place.Longitude))
	_, err = b.bot.Request(venueConfig)
	return err
}

func (b *Bot) handleLottery(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.LotteryData{}

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

func (b *Bot) handlePlayersList(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.PlayersListByGameData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	players, err := b.gamesFacade.GetPlayersByGameID(ctx, data.GameID)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		}

		return err
	}

	textBuilder := strings.Builder{}
	for i, player := range players {
		playerName := ""
		if player.UserID > 0 {
			var user model.User
			if user, err = b.usersFacade.GetUserByID(ctx, player.UserID); err != nil {
				return err
			}
			playerName = user.Name
		} else {
			var user model.User
			if user, err = b.usersFacade.GetUserByID(ctx, player.RegisteredBy); err != nil {
				return err
			}
			playerName = fmt.Sprintf("%s %s", getTranslator(legionerByLexeme)(ctx), user.Name)
		}

		if player.Degree == int32(registrator.Degree_DEGREE_UNLIKELY) {
			textBuilder.WriteString(fmt.Sprintf("%d. %s (%s)\n", i+1, playerName, getTranslator(degreeMap[registrator.Degree_DEGREE_UNLIKELY])(ctx)))
		} else {
			textBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, playerName))
		}
	}

	text := textBuilder.String()
	if text == "" {
		text = fmt.Sprintf("%s", getTranslator(listOfPlayersIsEmptyLexeme)(ctx))
	}

	msg := tgbotapi.NewMessage(clientID, text)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleRegisterGame(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.RegisterGameData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.gamesFacade.RegisterGame(ctx, data.GameID)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
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
	menu, err = b.getGameMenu(ctx, game, 0)
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

func (b *Bot) handleRegisterPlayer(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.RegisterPlayerData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.gamesFacade.RegisterPlayer(ctx, data.GameID, data.PlayerType, data.Degree)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		} else if errors.Is(err, model.ErrNoFreeSlot) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(noFreeSlotLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
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
	menu, err = b.getGameMenu(ctx, game, 0)
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

func (b *Bot) handleUnregisterGame(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.UnregisterGameData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.gamesFacade.UnregisterGame(ctx, data.GameID)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
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
	menu, err = b.getGameMenu(ctx, game, 0)
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

func (b *Bot) handleUnregisterPlayer(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.UnregisterPlayerData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	_, err = b.gamesFacade.UnregisterPlayer(ctx, data.GameID, data.PlayerType)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
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
	menu, err = b.getGameMenu(ctx, game, 0)
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

func (b *Bot) handleUpdatePayment(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.UpdatePaymentData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	err = b.gamesFacade.UpdatePayment(ctx, data.GameID, data.Payment)
	if err != nil {
		if errors.Is(err, model.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, getTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
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
	menu, err = b.getGameMenu(ctx, game, 1)
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

func (b *Bot) updateUserState(ctx context.Context, update *tgbotapi.Update, state int32) error {
	clientID := update.CallbackQuery.From.ID

	user, err := b.usersFacade.GetUserByTelegramID(ctx, clientID)
	if err != nil {
		return err
	}

	err = b.usersFacade.UpdateUserState(ctx, user.ID, state)
	if err != nil {
		return err
	}

	msg := tgbotapi.MessageConfig{}
	switch registrator.UserState(state) {
	case registrator.UserState_USER_STATE_CHANGING_EMAIL:
		msg = tgbotapi.NewMessage(clientID, getTranslator(enterYourEmailLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	case registrator.UserState_USER_STATE_CHANGING_NAME:
		msg = tgbotapi.NewMessage(clientID, getTranslator(enterYourNameLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	case registrator.UserState_USER_STATE_CHANGING_PHONE:
		msg = tgbotapi.NewMessage(clientID, getTranslator(enterYourPhoneLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
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
	registerStatus := fmt.Sprintf("%s %s", icons.UnregisteredGame, getTranslator(unregisteredGameLexeme)(ctx))
	if game.Registered {
		registerStatus = fmt.Sprintf("%s %s", icons.RegisteredGame, getTranslator(registeredGameLexeme)(ctx))
	}

	info.WriteString(registerStatus + "\n")

	paymentType := ""
	if strings.Index(game.PaymentType, "cash") != -1 {
		paymentType += strings.ToLower(getTranslator(cashLexeme)(ctx))
	}
	if strings.Index(game.PaymentType, "card") != -1 {
		if paymentType != "" {
			paymentType += ", "
		}
		paymentType += strings.ToLower(getTranslator(cardLexeme)(ctx))
	}

	if paymentType == "" {
		paymentType = "?"
	}

	if game.Payment != model.PaymentTypeInvalid {
		paymentStatus := fmt.Sprintf("%s %s: %s (%s)", icons.MixGamePayment, getTranslator(paymentLexeme)(ctx), strings.ToLower(getTranslator(mixLexeme)(ctx)), paymentType)
		if game.Payment == model.PaymentTypeCash {
			paymentStatus = fmt.Sprintf("%s %s: %s", icons.CashGamePayment, getTranslator(paymentLexeme)(ctx), paymentType)
		} else if game.Payment == model.PaymentTypeCertificate {
			paymentStatus = fmt.Sprintf("%s %s: %s", icons.FreeGamePayment, getTranslator(paymentLexeme)(ctx), strings.ToLower(getTranslator(certificateLexeme)(ctx)))
		}

		info.WriteString(paymentStatus + "\n")
	} else {
		info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.CashGamePayment, getTranslator(paymentLexeme)(ctx), paymentType))
	}

	if game.Name != "" {
		info.WriteString(fmt.Sprintf("%s %s: %s %s\n", icons.Sharp, getTranslator(titleLexeme)(ctx), game.Name, game.Number))
	} else {
		info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Sharp, getTranslator(numberLexeme)(ctx), game.Number))
	}

	if game.Price > 0 {
		price := strconv.Itoa(int(game.Price))
		info.WriteString(fmt.Sprintf("%s %s: %s₽\n", icons.USD, getTranslator(gameCostLexeme)(ctx), price))
	}

	info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Place, getTranslator(addressLexeme)(ctx), game.Place.Address))
	info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Calendar, getTranslator(dateTimeLexeme)(ctx), game.DateTime().String()))
	info.WriteString(fmt.Sprintf("%s %s: %d/%d/%d", icons.NumberOfPlayers, getTranslator(numberOfPlayersLexeme)(ctx), game.NumberOfPlayers, game.NumberOfLegioners, game.MaxPlayers))

	return info.String()
}
