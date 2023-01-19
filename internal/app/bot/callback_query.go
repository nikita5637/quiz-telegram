package bot

import (
	"context"
	"encoding/json"
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

	r, err := b.requestsFacade.GetRequest(ctx, update.CallbackData())
	if err != nil {
		return fmt.Errorf("get telegram request error: %w", err)
	}

	telegramRequest := TelegramRequest{}
	err = json.Unmarshal(r, &telegramRequest)
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
	case CommandGamesList:
		handler = func(ctx context.Context) error {
			return b.handleGamesList(ctx, update, telegramRequest)
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

	return b.unregisterRequest(ctx, update.CallbackData())
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

func (b *Bot) handleGamesList(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID

	resp, err := b.registratorServiceClient.GetGames(ctx, &registrator.GetGamesRequest{
		Active: true,
	})
	if err != nil {
		return err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, pbGame := range resp.GetGames() {
		leagueResp, err := b.registratorServiceClient.GetLeagueByID(ctx, &registrator.GetLeagueByIDRequest{
			Id: pbGame.GetLeagueId(),
		})
		if err != nil {
			return err
		}

		placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
			Id: pbGame.GetPlaceId(),
		})
		if err != nil {
			return err
		}

		pbReq := &registrator.GetGameByIDRequest{
			GameId: pbGame.GetId(),
		}

		var request model.Request
		request, err = getRequest(ctx, CommandGetGame, pbReq)
		if err != nil {
			return err
		}

		callbackData := b.registerRequest(ctx, request)
		text := fmt.Sprintf(gameInfoFormatString, leagueResp.GetLeague().GetShortName(), pbGame.GetNumber(), placeResp.GetPlace().GetShortName(), model.DateTime(pbGame.GetDate().AsTime()))

		if pbGame.GetMy() {
			text = myGamePrefix + text
		} else {
			if pbGame.GetNumberOfLegioners()+pbGame.GetNumberOfPlayers() > 0 {
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

	req := &registrator.GetGameByIDRequest{}

	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	resp, err := b.registratorServiceClient.GetGameByID(ctx, req)
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

	game := convertPBGameToModelGame(resp.GetGame())

	game.Address = "?"

	placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: resp.GetGame().GetPlaceId(),
	})
	if err != nil {
		logger.Warnf(ctx, "getting place info error: %w", err)
	} else {
		game.Address = placeResp.GetPlace().GetAddress()
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

	req := &registrator.GetPhotosByGameIDRequest{}
	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	resp, err := b.photographerServiceClient.GetPhotosByGameID(ctx, req)
	if err != nil {
		st := status.Convert(err)
		// unlikely
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

	deleteConfig := tgbotapi.NewDeleteMessage(clientID, messageID)
	_, err = b.bot.Request(deleteConfig)
	if err != nil {
		return err
	}

	for _, url := range resp.GetUrls() {
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

	req := &registrator.GetGamesWithPhotosRequest{}

	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	resp, err := b.photographerServiceClient.GetGamesWithPhotos(ctx, req)
	if err != nil {
		return err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, pbGame := range resp.GetGames() {
		leagueResp, err := b.registratorServiceClient.GetLeagueByID(ctx, &registrator.GetLeagueByIDRequest{
			Id: pbGame.GetLeagueId(),
		})
		if err != nil {
			return err
		}

		placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
			Id: pbGame.GetPlaceId(),
		})
		if err != nil {
			return err
		}

		pbReq := &registrator.GetPhotosByGameIDRequest{
			GameId: pbGame.GetId(),
		}

		var request model.Request
		request, err = getRequest(ctx, CommandGetGamePhotos, pbReq)
		if err != nil {
			return err
		}

		callbackData := b.registerRequest(ctx, request)
		text := fmt.Sprintf(gamePhotosInfoFormatString, model.ResultPlace(pbGame.GetResultPlace()).String(), leagueResp.GetLeague().GetShortName(), pbGame.GetNumber(), placeResp.GetPlace().GetShortName(), model.DateTime(pbGame.GetDate().AsTime()))

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	navigateButtonsRow := make([]tgbotapi.InlineKeyboardButton, 0, 2)

	if req.GetOffset() > 0 {
		offset := uint32(0)
		if req.GetOffset() > gamesWithPhotosListLimit {
			offset = req.GetOffset() - gamesWithPhotosListLimit
		}

		pbReq := &registrator.GetGamesWithPhotosRequest{
			Limit:  gamesWithPhotosListLimit,
			Offset: offset,
		}

		var request model.Request
		request, err = getRequest(ctx, CommandGetListGamesWithPhotosPrevPage, pbReq)
		if err != nil {
			return err
		}

		callbackData := b.registerRequest(ctx, request)

		btnPrev := tgbotapi.InlineKeyboardButton{
			Text:         prevPageStringText,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnPrev)
	}

	leftNext := uint32(0)
	if resp.GetTotal() > (req.GetOffset() + req.GetLimit()) {
		leftNext = resp.GetTotal() - (req.GetOffset() + req.GetLimit())
	}

	if leftNext > 0 {
		pbReq := &registrator.GetGamesWithPhotosRequest{
			Limit:  gamesWithPhotosListLimit,
			Offset: req.GetOffset() + req.GetLimit(),
		}

		var request model.Request
		request, err = getRequest(ctx, CommandGetListGamesWithPhotosNextPage, pbReq)
		if err != nil {
			return err
		}

		callbackData := b.registerRequest(ctx, request)

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

	req := &registrator.GetGamesWithPhotosRequest{}

	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	resp, err := b.photographerServiceClient.GetGamesWithPhotos(ctx, req)
	if err != nil {
		return err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, pbGame := range resp.GetGames() {
		leagueResp, err := b.registratorServiceClient.GetLeagueByID(ctx, &registrator.GetLeagueByIDRequest{
			Id: pbGame.GetLeagueId(),
		})
		if err != nil {
			return err
		}

		placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
			Id: pbGame.GetPlaceId(),
		})
		if err != nil {
			return err
		}

		pbReq := &registrator.GetPhotosByGameIDRequest{
			GameId: pbGame.GetId(),
		}

		var request model.Request
		request, err = getRequest(ctx, CommandGetGamePhotos, pbReq)
		if err != nil {
			return err
		}

		callbackData := b.registerRequest(ctx, request)
		text := fmt.Sprintf(gamePhotosInfoFormatString, model.ResultPlace(pbGame.GetResultPlace()).String(), leagueResp.GetLeague().GetShortName(), pbGame.GetNumber(), placeResp.GetPlace().GetShortName(), model.DateTime(pbGame.GetDate().AsTime()))

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	navigateButtonsRow := make([]tgbotapi.InlineKeyboardButton, 0, 2)

	if req.GetOffset() > 0 {
		offset := uint32(0)
		if req.GetOffset() > gamesWithPhotosListLimit {
			offset = req.GetOffset() - gamesWithPhotosListLimit
		}

		pbReq := &registrator.GetGamesWithPhotosRequest{
			Limit:  gamesWithPhotosListLimit,
			Offset: offset,
		}

		var request model.Request
		request, err = getRequest(ctx, CommandGetListGamesWithPhotosPrevPage, pbReq)
		if err != nil {
			return err
		}

		callbackData := b.registerRequest(ctx, request)

		btnPrev := tgbotapi.InlineKeyboardButton{
			Text:         prevPageStringText,
			CallbackData: &callbackData,
		}
		navigateButtonsRow = append(navigateButtonsRow, btnPrev)
	}

	leftNext := uint32(0)
	if resp.GetTotal() > (req.GetOffset() + req.GetLimit()) {
		leftNext = resp.GetTotal() - (req.GetOffset() + req.GetLimit())
	}

	if leftNext > 0 {
		pbReq := &registrator.GetGamesWithPhotosRequest{
			Limit:  gamesWithPhotosListLimit,
			Offset: req.GetOffset() + req.GetLimit(),
		}

		var request model.Request
		request, err = getRequest(ctx, CommandGetListGamesWithPhotosNextPage, pbReq)
		if err != nil {
			return err
		}

		callbackData := b.registerRequest(ctx, request)

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

func (b *Bot) handleLottery(ctx context.Context, update *tgbotapi.Update, telegramRequest TelegramRequest) error {
	clientID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	req := &registrator.RegisterForLotteryRequest{}

	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	resp, err := b.croupierServiceClient.RegisterForLottery(ctx, req)
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

	req := &registrator.GetPlayersByGameIDRequest{}
	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	resp, err := b.registratorServiceClient.GetPlayersByGameID(ctx, req)
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
			if playerResp, err := b.registratorServiceClient.GetUserByID(ctx, &registrator.GetUserByIDRequest{
				Id: player.UserId,
			}); err != nil {
				return err
			} else {
				playerName = playerResp.GetUser().GetName()
			}
		} else {
			if playerResp, err := b.registratorServiceClient.GetUserByID(ctx, &registrator.GetUserByIDRequest{
				Id: player.RegisteredBy,
			}); err != nil {
				return err
			} else {
				playerName = fmt.Sprintf("%s %s", getTranslator(legionerByLexeme)(ctx), playerResp.GetUser().GetName())
			}
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

	req := &registrator.RegisterGameRequest{}

	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.RegisterGame(ctx, req)
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

	gameResp, err := b.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
		GameId: req.GetGameId(),
	})
	if err != nil {
		return err
	}

	game := convertPBGameToModelGame(gameResp.GetGame())

	game.Address = "?"

	placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: gameResp.GetGame().GetPlaceId(),
	})
	if err != nil {
		logger.Warnf(ctx, "getting place info error: %w", err)
	} else {
		game.Address = placeResp.GetPlace().GetAddress()
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

	req := &registrator.RegisterPlayerRequest{}
	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.RegisterPlayer(ctx, req)
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

	gameResp, err := b.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
		GameId: req.GetGameId(),
	})
	if err != nil {
		return err
	}

	game := convertPBGameToModelGame(gameResp.GetGame())

	game.Address = "?"

	placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: gameResp.GetGame().GetPlaceId(),
	})
	if err != nil {
		logger.Warnf(ctx, "getting place info error: %w", err)
	} else {
		game.Address = placeResp.GetPlace().GetAddress()
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

	req := &registrator.UnregisterGameRequest{}
	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.UnregisterGame(ctx, req)
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

	gameResp, err := b.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
		GameId: req.GetGameId(),
	})
	if err != nil {
		return err
	}

	game := convertPBGameToModelGame(gameResp.GetGame())

	game.Address = "?"

	placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: gameResp.GetGame().GetPlaceId(),
	})
	if err != nil {
		logger.Warnf(ctx, "getting place info error: %w", err)
	} else {
		game.Address = placeResp.GetPlace().GetAddress()
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

	req := &registrator.UnregisterPlayerRequest{}
	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.UnregisterPlayer(ctx, req)
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

	gameResp, err := b.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
		GameId: req.GetGameId(),
	})
	if err != nil {
		return err
	}

	game := convertPBGameToModelGame(gameResp.GetGame())

	game.Address = "?"

	placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: gameResp.GetGame().GetPlaceId(),
	})
	if err != nil {
		logger.Warnf(ctx, "getting place info error: %w", err)
	} else {
		game.Address = placeResp.GetPlace().GetAddress()
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

	req := &registrator.UpdatePaymentRequest{}
	err := json.Unmarshal(telegramRequest.Body, req)
	if err != nil {
		return err
	}

	_, err = b.registratorServiceClient.UpdatePayment(ctx, req)
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

	gameResp, err := b.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
		GameId: req.GetGameId(),
	})
	if err != nil {
		return err
	}

	game := convertPBGameToModelGame(gameResp.GetGame())

	game.Address = "?"

	placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
		Id: gameResp.GetGame().GetPlaceId(),
	})
	if err != nil {
		logger.Warnf(ctx, "getting place info error: %w", err)
	} else {
		game.Address = placeResp.GetPlace().GetAddress()
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
	info.WriteString("📍 Адрес: " + game.Address + "\n")
	info.WriteString("📅 Дата и время: " + game.DateTime().String() + "\n")
	info.WriteString(fmt.Sprintf("👥 Количество игроков: %d/%d/%d", game.NumberPlayers, game.NumberLegioners, game.MaxPlayers))

	return info.String()
}
