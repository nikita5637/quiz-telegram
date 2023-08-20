package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mono83/maybe"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gameplayers"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	callbackdata_utils "github.com/nikita5637/quiz-telegram/internal/pkg/utils/callbackdata"
	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"
	user_utils "github.com/nikita5637/quiz-telegram/utils/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	croupierpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/croupier"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
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
	degreeMap = map[model.Degree]i18n.Lexeme{
		model.DegreeInvalid: {
			Key:      "invalid_degree",
			FallBack: "invalid degree",
		},
		model.DegreeLikely: {
			Key:      "plays_likely",
			FallBack: "plays likely",
		},
		model.DegreeUnlikely: {
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
	enterYourBirthdateLexeme = i18n.Lexeme{
		Key:      "enter_your_birthdate",
		FallBack: "OK. Enter your birthdate(format DD.MM.YYYY).",
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
	enterYourSexLexeme = i18n.Lexeme{
		Key:      "enter_your_sex",
		FallBack: "OK. Enter your sex",
	}
	gameCostLexeme = i18n.Lexeme{
		Key:      "game_cost",
		FallBack: "Game cost",
	}
	gameHasPassedLexeme = i18n.Lexeme{
		Key:      "game_has_passed",
		FallBack: "Game has passed",
	}
	gameNotFoundLexeme = i18n.Lexeme{
		Key:      "game_not_found",
		FallBack: "Game not found",
	}
	legionerByLexeme = i18n.Lexeme{
		Key:      "legioner_by",
		FallBack: "Legioner by",
	}
	legionerIsSignedUpForTheGameLexeme = i18n.Lexeme{
		Key:      "legioner_is_signed_up_for_the_game",
		FallBack: "Legioner is signed up for the game",
	}
	legionerIsSignedUpForTheGameUnlikelyLexeme = i18n.Lexeme{
		Key:      "legioner_is_signed_up_for_the_game_unlikely",
		FallBack: "Legioner is signed up for the game unlikely",
	}
	legionerIsUnsignedUpForTheGameLexeme = i18n.Lexeme{
		Key:      "legioner_is_unsigned_up_for_the_game",
		FallBack: "Legioner is unsigned up for the game",
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
	thereAreNoYourLegionersRegisteredForTheGameLexeme = i18n.Lexeme{
		Key:      "there_are_no_your_legioners_registered_for_the_game",
		FallBack: "There are no your legioners registered for the game",
	}
	titleLexeme = i18n.Lexeme{
		Key:      "title",
		FallBack: "Title",
	}
	unregisteredGameLexeme = i18n.Lexeme{
		Key:      "unregistered_game",
		FallBack: "We are unregistered for the game",
	}
	youAreNotRegisteredForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_not_registered_for_the_game",
		FallBack: "You are not registered for the game",
	}
	youAreSignedUpForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_signed_up_for_the_game",
		FallBack: "You are signed up for the game",
	}
	youAreSignedUpForTheGameUnlikelyLexeme = i18n.Lexeme{
		Key:      "you_are_signed_up_for_the_game_unlikely",
		FallBack: "You are signed up for the game unlikely",
	}
	youAreUnsignedUpForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_unsigned_up_for_the_game",
		FallBack: "You are unsigned up for the game",
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

	user, err := b.checkAuth(ctx, clientID)
	if err != nil {
		name := update.CallbackQuery.Message.Chat.FirstName
		if name == "" {
			name = update.CallbackQuery.Message.Chat.UserName
		}

		_, err = b.usersFacade.CreateUser(ctx, name, clientID, int32(usermanagerpb.UserState_USER_STATE_WELCOME))
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

	ctx = user_utils.NewContextWithUser(ctx, user)

	telegramRequest := commands.TelegramRequest{}

	err = json.Unmarshal([]byte(update.CallbackData()), &telegramRequest)
	if err != nil {
		return fmt.Errorf("telegram request unmarshaling error: %w", err)
	}

	type handlerFunc func(ctx context.Context) error
	var handler handlerFunc
	switch telegramRequest.Command {
	case commands.CommandChangeBirthdate:
		handler = func(ctx context.Context) error {
			return b.handleChangeBirthdate(ctx, update, telegramRequest)
		}
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
	case commands.CommandChangeSex:
		handler = func(ctx context.Context) error {
			return b.handleChangeSex(ctx, update, telegramRequest)
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
			data := &commands.PlayersListByGameData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return fmt.Errorf("telegram request body unmarshal error: %w", err)
			}

			return b.handlePlayersList(ctx, update, data)
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
			data := &commands.UnregisterPlayerData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return fmt.Errorf("telegram request body unmarshal error: %w", err)
			}

			return b.handleUnregisterPlayer(ctx, update, data)
		}
	case commands.CommandUpdatePayment:
		handler = func(ctx context.Context) error {
			return b.handleUpdatePayment(ctx, update, telegramRequest)
		}
	case commands.CommandUpdatePlayerRegistration:
		handler = func(ctx context.Context) error {
			data := &commands.UpdatePlayerRegistration{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return fmt.Errorf("telegram request body unmarshal error: %w", err)
			}

			return b.handleUpdatePlayerRegistration(ctx, update, data)
		}
	}

	err = handler(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleChangeBirthdate(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_BIRTHDATE))
}

func (b *Bot) handleChangeEmail(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_EMAIL))
}

func (b *Bot) handleChangeName(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME))
}

func (b *Bot) handleChangePhone(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_PHONE))
}

func (b *Bot) handleChangeSex(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_SEX))
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
		callbackData, err = callbackdata_utils.GetCallbackData(ctx, commands.CommandGetGame, payload)
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

	msg := tgbotapi.NewEditMessageText(clientID, update.CallbackQuery.Message.MessageID, i18n.GetTranslator(chooseGameLexeme)(ctx))
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
		if errors.Is(err, games.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		}

		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
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

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
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
		if errors.Is(err, games.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
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
		_, err = b.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
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
		callbackData, err = callbackdata_utils.GetCallbackData(ctx, commands.CommandGetGamePhotos, payload)
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
		callbackData, err = callbackdata_utils.GetCallbackData(ctx, commands.CommandGetListGamesWithPhotosPrevPage, payload)
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
		callbackData, err = callbackdata_utils.GetCallbackData(ctx, commands.CommandGetListGamesWithPhotosNextPage, payload)
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
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
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
		callbackData, err = callbackdata_utils.GetCallbackData(ctx, commands.CommandGetGamePhotos, payload)
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
		callbackData, err = callbackdata_utils.GetCallbackData(ctx, commands.CommandGetListGamesWithPhotosPrevPage, payload)
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
		callbackData, err = callbackdata_utils.GetCallbackData(ctx, commands.CommandGetListGamesWithPhotosNextPage, payload)
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
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
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
	if err != nil {
		logger.Errorf(ctx, "sending venue error: %s", err)
		return err
	}

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleLottery(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	data := &commands.LotteryData{}

	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	resp, err := b.croupierServiceClient.RegisterForLottery(ctx, &croupierpb.RegisterForLotteryRequest{
		GameId: data.GameID,
	})
	if err != nil {
		st := status.Convert(err)

		msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(somethingWentWrongLexeme)(ctx))
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

	msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(youHaveSuccessfullyRegisteredInLotteryLexeme)(ctx))
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	if resp.GetNumber() > 0 {
		var newMsg tgbotapi.Message
		msg := tgbotapi.NewMessage(clientID, fmt.Sprintf("%s: %d", i18n.GetTranslator(yourLotteryNumberIsLexeme)(ctx), resp.GetNumber()))
		newMsg, err = b.bot.Send(msg)
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

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handlePlayersList(ctx context.Context, update *tgbotapi.Update, data *commands.PlayersListByGameData) error {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.PlayersListByGameData) (*tgbotapi.MessageConfig, *tgbotapi.CallbackConfig, error) {
		clientID := update.CallbackQuery.From.ID
		gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("get game players by game ID error: %w", err)
		}

		if len(gamePlayers) == 0 {
			msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfPlayersIsEmptyLexeme)(ctx))
			return &msg, nil, nil
		}

		textBuilder := strings.Builder{}
		for i, gamePlayer := range gamePlayers {
			playerName := ""
			if userID, ok := gamePlayer.UserID.Get(); ok {
				var user model.User
				if user, err = b.usersFacade.GetUserByID(ctx, userID); err != nil {
					return nil, nil, fmt.Errorf("get user by ID error: %w", err)
				}
				playerName = user.Name
			} else {
				var user model.User
				if user, err = b.usersFacade.GetUserByID(ctx, gamePlayer.RegisteredBy); err != nil {
					return nil, nil, fmt.Errorf("get user by ID error: %w", err)
				}
				playerName = fmt.Sprintf("%s %s", i18n.GetTranslator(legionerByLexeme)(ctx), user.Name)
			}

			if gamePlayer.Degree == model.DegreeUnlikely {
				textBuilder.WriteString(fmt.Sprintf("%d. %s (%s)\n", i+1, playerName, i18n.GetTranslator(degreeMap[model.DegreeUnlikely])(ctx)))
			} else {
				textBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, playerName))
			}
		}

		msg := tgbotapi.NewMessage(clientID, textBuilder.String())
		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		return &msg, &cb, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return fmt.Errorf("prepare game players list message error: %w", err)
	}

	if msg != nil {
		if _, err := b.bot.Send(msg); err != nil {
			logger.Errorf(ctx, "sending message error: %s", err)
		}
	}

	if cb != nil {
		if _, err := b.bot.Request(cb); err != nil {
			logger.Errorf(ctx, "sending callback error: %s", err)
		}
	}

	return nil
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
		if errors.Is(err, games.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		}

		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
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

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(registeredGameLexeme)(ctx))
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleRegisterPlayer(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID
	callbackID := update.CallbackQuery.ID

	data := &commands.RegisterPlayerData{}
	err := json.Unmarshal(telegramRequest.Body, data)
	if err != nil {
		return err
	}

	gamePlayer := model.GamePlayer{
		GameID:       data.GameID,
		UserID:       maybe.Nothing[int32](),
		RegisteredBy: data.RegisteredBy,
		Degree:       data.Degree,
	}

	if data.UserID != 0 {
		gamePlayer.UserID = maybe.Just(data.UserID)
	}

	err = b.gamePlayersFacade.RegisterPlayer(ctx, gamePlayer)
	if err != nil {
		if errors.Is(err, games.ErrGameHasPassed) {
			msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(gameHasPassedLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		} else if errors.Is(err, games.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		} else if errors.Is(err, gameplayers.ErrNoFreeSlot) {
			msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(noFreeSlotLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		} else if errors.Is(err, gameplayers.ErrGamePlayerAlreadyRegistered) {
			msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(youAreAlreadyRegisteredForTheGameLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		}

		return err
	}

	cb := tgbotapi.NewCallback(callbackID, "")
	if data.UserID == 0 && data.RegisteredBy != data.UserID {
		if data.Degree == model.DegreeLikely {
			cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(legionerIsSignedUpForTheGameLexeme)(ctx))
		} else {
			cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(legionerIsSignedUpForTheGameUnlikelyLexeme)(ctx))
		}
	} else if data.UserID != 0 && data.RegisteredBy == data.UserID {
		if data.Degree == model.DegreeLikely {
			cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(youAreSignedUpForTheGameLexeme)(ctx))
		} else {
			cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(youAreSignedUpForTheGameUnlikelyLexeme)(ctx))
		}
	}

	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
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
		if errors.Is(err, games.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		}

		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
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

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(unregisteredGameLexeme)(ctx))
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleUnregisterPlayer(ctx context.Context, update *tgbotapi.Update, data *commands.UnregisterPlayerData) error {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.UnregisterPlayerData) (*tgbotapi.EditMessageTextConfig, *tgbotapi.CallbackConfig, error) {
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID
		callbackID := update.CallbackQuery.ID

		gamePlayer := model.GamePlayer{
			GameID:       data.GameID,
			UserID:       maybe.Nothing[int32](),
			RegisteredBy: data.RegisteredBy,
			Degree:       data.Degree,
		}
		if data.UserID != 0 {
			gamePlayer.UserID = maybe.Just(data.UserID)
		}

		if err := b.gamePlayersFacade.UnregisterPlayer(ctx, gamePlayer); err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
				return &msg, nil, nil
			} else if errors.Is(err, games.ErrGameHasPassed) {
				cb := tgbotapi.NewCallback(callbackID, i18n.GetTranslator(gameHasPassedLexeme)(ctx))
				return nil, &cb, nil
			} else if errors.Is(err, gameplayers.ErrGamePlayerNotFound) {
				var cb tgbotapi.CallbackConfig
				if data.UserID != 0 && data.RegisteredBy == data.UserID {
					cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(youAreNotRegisteredForTheGameLexeme)(ctx))
				} else {
					cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(thereAreNoYourLegionersRegisteredForTheGameLexeme)(ctx))
				}
				return nil, &cb, nil
			}

			return nil, nil, err
		}

		cb := tgbotapi.NewCallback(callbackID, "")
		if data.UserID == 0 && data.RegisteredBy != data.UserID {
			cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(legionerIsUnsignedUpForTheGameLexeme)(ctx))
		} else if data.UserID != 0 && data.RegisteredBy == data.UserID {
			cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(youAreUnsignedUpForTheGameLexeme)(ctx))
		}

		game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("get game by ID error: %w", err)
		}

		lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
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
			return nil, nil, fmt.Errorf("get game menu error: %w", err)
		}

		msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
		msg.ReplyMarkup = &menu

		return &msg, &cb, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return fmt.Errorf("prepare unregister player message error: %w", err)
	}

	if msg != nil {
		if _, err := b.bot.Send(msg); err != nil {
			logger.Errorf(ctx, "sending message error: %s", err)
		}
	}

	if cb != nil {
		if _, err := b.bot.Request(cb); err != nil {
			logger.Errorf(ctx, "sending callback error: %s", err)
		}
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
		if errors.Is(err, games.ErrGameNotFound) {
			msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
			_, err = b.bot.Send(msg)
			return err
		}

		return err
	}

	game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
	if err != nil {
		return err
	}

	lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
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

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	switch game.Payment {
	case model.PaymentTypeCash:
		cb = tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(cashGamePaymentLexeme)(ctx))
	case model.PaymentTypeCertificate:
		cb = tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(freeGamePaymentLexeme)(ctx))
	case model.PaymentTypeMixed:
		cb = tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(mixGamePaymentLexeme)(ctx))
	}

	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
}

func (b *Bot) handleUpdatePlayerRegistration(ctx context.Context, update *tgbotapi.Update, data *commands.UpdatePlayerRegistration) error {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.UpdatePlayerRegistration) (*tgbotapi.EditMessageTextConfig, *tgbotapi.CallbackConfig, error) {

		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID
		callbackID := update.CallbackQuery.ID

		gamePlayer := model.GamePlayer{
			GameID:       data.GameID,
			UserID:       maybe.Just(data.UserID),
			RegisteredBy: data.RegisteredBy,
			Degree:       data.Degree,
		}

		if err := b.gamePlayersFacade.UpdatePlayerRegistration(ctx, gamePlayer); err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(gameNotFoundLexeme)(ctx))
				return &msg, nil, nil
			}

			return nil, nil, fmt.Errorf("update player registration error: %w", err)
		}

		cb := tgbotapi.NewCallback(callbackID, "")
		if data.UserID == 0 && data.RegisteredBy != data.UserID {
			if data.Degree == model.DegreeUnlikely {
				cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(legionerIsSignedUpForTheGameUnlikelyLexeme)(ctx))
			} else {
				cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(legionerIsSignedUpForTheGameLexeme)(ctx))
			}
		} else if data.UserID != 0 && data.RegisteredBy == data.UserID {
			if data.Degree == model.DegreeUnlikely {
				cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(youAreSignedUpForTheGameUnlikelyLexeme)(ctx))
			} else {
				cb = tgbotapi.NewCallback(callbackID, i18n.GetTranslator(youAreSignedUpForTheGameLexeme)(ctx))
			}
		}

		game, err := b.gamesFacade.GetGameByID(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("get game by ID error: %w", err)
		}

		lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
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
			return nil, nil, fmt.Errorf("get game menu error: %w", err)
		}

		msg := tgbotapi.NewEditMessageText(clientID, messageID, detailInfo(ctx, game))
		msg.ReplyMarkup = &menu

		return &msg, &cb, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return fmt.Errorf("prepare update game player registration message error: %w", err)
	}

	if cb != nil {
		if _, err := b.bot.Request(cb); err != nil {
			logger.Errorf(ctx, "sending callback error: %s", err)
		}
	}

	if msg != nil {
		if _, err := b.bot.Send(msg); err != nil {
			logger.Errorf(ctx, "sending message error: %s", err)
		}
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
	switch usermanagerpb.UserState(state) {
	case usermanagerpb.UserState_USER_STATE_CHANGING_BIRTHDATE:
		msg = tgbotapi.NewMessage(clientID, i18n.GetTranslator(enterYourBirthdateLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	case usermanagerpb.UserState_USER_STATE_CHANGING_EMAIL:
		msg = tgbotapi.NewMessage(clientID, i18n.GetTranslator(enterYourEmailLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	case usermanagerpb.UserState_USER_STATE_CHANGING_NAME:
		msg = tgbotapi.NewMessage(clientID, i18n.GetTranslator(enterYourNameLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	case usermanagerpb.UserState_USER_STATE_CHANGING_PHONE:
		msg = tgbotapi.NewMessage(clientID, i18n.GetTranslator(enterYourPhoneLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	case usermanagerpb.UserState_USER_STATE_CHANGING_SEX:
		msg = tgbotapi.NewMessage(clientID, i18n.GetTranslator(enterYourSexLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	}

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf(ctx, "sending message error: %s", err)
		return err
	}

	cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	_, err = b.bot.Request(cb)
	if err != nil {
		logger.Errorf(ctx, "sending callback error: %s", err)
		return err
	}

	return nil
}

func detailInfo(ctx context.Context, game model.Game) string {
	info := strings.Builder{}
	registerStatus := fmt.Sprintf("%s %s", icons.UnregisteredGame, i18n.GetTranslator(unregisteredGameLexeme)(ctx))
	if game.Registered {
		registerStatus = fmt.Sprintf("%s %s", icons.RegisteredGame, i18n.GetTranslator(registeredGameLexeme)(ctx))
	}

	info.WriteString(registerStatus + "\n")

	paymentType := ""
	if strings.Index(game.PaymentType, "cash") != -1 {
		paymentType += strings.ToLower(i18n.GetTranslator(cashLexeme)(ctx))
	}
	if strings.Index(game.PaymentType, "card") != -1 {
		if paymentType != "" {
			paymentType += ", "
		}
		paymentType += strings.ToLower(i18n.GetTranslator(cardLexeme)(ctx))
	}

	if paymentType == "" {
		paymentType = "?"
	}

	if game.Payment != model.PaymentTypeInvalid {
		paymentStatus := fmt.Sprintf("%s %s: %s (%s)", icons.MixGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), strings.ToLower(i18n.GetTranslator(mixLexeme)(ctx)), paymentType)
		if game.Payment == model.PaymentTypeCash {
			paymentStatus = fmt.Sprintf("%s %s: %s", icons.CashGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), paymentType)
		} else if game.Payment == model.PaymentTypeCertificate {
			paymentStatus = fmt.Sprintf("%s %s: %s", icons.FreeGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), strings.ToLower(i18n.GetTranslator(certificateLexeme)(ctx)))
		}

		info.WriteString(paymentStatus + "\n")
	} else {
		info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.CashGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), paymentType))
	}

	if game.Name != "" {
		info.WriteString(fmt.Sprintf("%s %s: %s %s\n", icons.Sharp, i18n.GetTranslator(titleLexeme)(ctx), game.Name, game.Number))
	} else {
		info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Sharp, i18n.GetTranslator(numberLexeme)(ctx), game.Number))
	}

	if game.Price > 0 {
		price := strconv.Itoa(int(game.Price))
		info.WriteString(fmt.Sprintf("%s %s: %s₽\n", icons.USD, i18n.GetTranslator(gameCostLexeme)(ctx), price))
	}

	info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Place, i18n.GetTranslator(addressLexeme)(ctx), game.Place.Address))
	info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Calendar, i18n.GetTranslator(dateTimeLexeme)(ctx), game.DateTime().String()))
	info.WriteString(fmt.Sprintf("%s %s: %d/%d/%d", icons.NumberOfPlayers, i18n.GetTranslator(numberOfPlayersLexeme)(ctx), game.NumberOfPlayers, game.NumberOfLegioners, game.MaxPlayers))

	return info.String()
}
