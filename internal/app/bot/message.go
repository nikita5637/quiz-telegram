package bot

import (
	"context"
	"fmt"

	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
)

const (
	gameInfoFormatString       = "%s %s %s %s\n"
	gamePhotosInfoFormatString = "📸" + "%s" + gameInfoFormatString
	gameWithPlayersPrefix      = "❗️ "
	nextPageStringText         = ">"
	myGamePrefix               = "ℹ️ "
	settingFormatString        = "%s [%s]"
)

var (
	changeEmailLexeme = i18n.Lexeme{
		Key:      "change_email",
		FallBack: "Change email",
	}
	changeNameLexeme = i18n.Lexeme{
		Key:      "change_name",
		FallBack: "Change name",
	}
	changePhoneLexeme = i18n.Lexeme{
		Key:      "change_phone",
		FallBack: "Change phone",
	}
	chooseGameLexeme = i18n.Lexeme{
		Key:      "choose_a_game",
		FallBack: "Choose a game",
	}
	emailChangedLexeme = i18n.Lexeme{
		Key:      "email_changed",
		FallBack: "Email changed",
	}
	gamePhotosLexeme = i18n.Lexeme{
		Key:      "game_photos",
		FallBack: "Game photos",
	}
	helpMessageLexeme = i18n.Lexeme{
		Key:      "help_message",
		FallBack: "Help message",
	}
	listOfGamesIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_games_is_empty",
		FallBack: "There are not games",
	}
	listOfGamesWithPhotosIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_games_with_photos_is_empty",
		FallBack: "There are not games with photos",
	}
	listOfMyGamesIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_my_games_is_empty",
		FallBack: "You don't play with us yet",
	}
	listOfMyGamesLexeme = i18n.Lexeme{
		Key:      "list_of_your_games",
		FallBack: "List of your games",
	}
	listOfRegisteredGamesLexeme = i18n.Lexeme{
		Key:      "list_of_registered_games",
		FallBack: "List of registered games",
	}
	listOfRegisteredGamesIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_registered_games_is_empty",
		FallBack: "There are not registered games",
	}
	nameChangedLexeme = i18n.Lexeme{
		Key:      "name_changed",
		FallBack: "Name changed",
	}
	permissionDeniedLexeme = i18n.Lexeme{
		Key:      "permission_denied",
		FallBack: "Permission denied",
	}
	phoneChangedLexeme = i18n.Lexeme{
		Key:      "phone_changed",
		FallBack: "Phone changed",
	}
	settingsLexeme = i18n.Lexeme{
		Key:      "settings",
		FallBack: "Settings",
	}
	somethingWentWrongLexeme = i18n.Lexeme{
		Key:      "something_went_wrong",
		FallBack: "Something went wrong",
	}
	welcomeMessageLexeme = i18n.Lexeme{
		Key:      "welcome_message",
		FallBack: "Welcome message",
	}
)

// HandleMessage ...
func (b *Bot) HandleMessage(ctx context.Context, update *tgbotapi.Update) error {
	if update.Message.Chat.IsSuperGroup() {
		logger.DebugKV(ctx, "skipped supergroup message", "groupID", update.Message.Chat.ID)
		return nil
	}

	clientID := update.Message.From.ID
	firstName := update.Message.From.FirstName
	userName := update.Message.From.UserName
	text := update.Message.Text

	ctx = telegram_utils.NewContextWithClientID(ctx, clientID)

	logger.DebugKV(ctx, "new private message incoming", "clientID", clientID, "text", text)

	err := b.checkAuth(ctx, clientID)
	if err != nil {
		name := firstName
		if name == "" {
			name = userName
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

	var handler func(ctx context.Context) error

	switch text {
	case "/games":
		handler = func(ctx context.Context) error {
			var msg tgbotapi.Chattable
			msg, err = b.getListOfGamesMessage(ctx, update)
			if err != nil {
				return err
			}

			_, err = b.bot.Send(msg)
			if err != nil {
				logger.Errorf(ctx, "error while sending message: %w", err)
			}

			return err
		}
	case "/help":
		handler = func(ctx context.Context) error {
			helpMessage := helpMessage(ctx, clientID)
			_, err = b.bot.Send(helpMessage)
			return err
		}
	case "/mygames":
		handler = func(ctx context.Context) error {
			var msg tgbotapi.Chattable
			msg, err = b.getListOfMyGamesMessage(ctx, update)
			if err != nil {
				return err
			}

			_, err = b.bot.Send(msg)
			if err != nil {
				logger.Errorf(ctx, "error while sending message: %w", err)
			}

			return err
		}
	case "/photos":
		handler = func(ctx context.Context) error {
			var msg tgbotapi.Chattable
			msg, err = b.getGamesWithPhotosMessage(ctx, update)
			if err != nil {
				return err
			}

			_, err = b.bot.Send(msg)
			if err != nil {
				logger.Errorf(ctx, "error while sending message: %w", err)
			}

			return err
		}
	case "/registeredgames":
		handler = func(ctx context.Context) error {
			var msg tgbotapi.Chattable
			msg, err = b.getListOfRegisteredGamesMessage(ctx, update)
			if err != nil {
				return err
			}

			_, err = b.bot.Send(msg)
			if err != nil {
				logger.Errorf(ctx, "error while sending message: %w", err)
			}

			return err
		}
	case "/settings":
		handler = func(ctx context.Context) error {
			var settingsMessage tgbotapi.Chattable
			settingsMessage, err = b.getSettingsMessage(ctx, update)
			if err != nil {
				return err
			}

			_, err = b.bot.Send(settingsMessage)
			if err != nil {
				logger.Errorf(ctx, "error while sending message: %w", err)
			}

			return err
		}
	default:
		handler = func(ctx context.Context) error {
			return b.handleDefaultMessage(ctx, update)
		}
	}

	if err := handler(ctx); err != nil {
		responseMessage := tgbotapi.NewMessage(clientID, getTranslator(somethingWentWrongLexeme)(ctx))
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.PermissionDenied {
				for _, detail := range st.Details() {
					switch t := detail.(type) {
					case *errdetails.ErrorInfo:
						reason := t.GetReason()
						if reason == "banned" {
							responseMessage = tgbotapi.NewMessage(clientID, getTranslator(permissionDeniedLexeme)(ctx))
						}
					}
				}
			} else if st.Code() == codes.InvalidArgument {
				for _, detail := range st.Details() {
					switch t := detail.(type) {
					case *errdetails.LocalizedMessage:
						localizedMessage := t.GetMessage()
						responseMessage = tgbotapi.NewMessage(clientID, localizedMessage)
					}
				}
			}
		}

		_, err = b.bot.Send(responseMessage)
		if err != nil {
			logger.Errorf(ctx, "error while send message: %s", err.Error())
		}
	}

	return nil
}

func (b *Bot) handleDefaultMessage(ctx context.Context, update *tgbotapi.Update) error {
	clientID := update.Message.From.ID
	resp, err := b.registratorServiceClient.GetUserByTelegramID(ctx, &registrator.GetUserByTelegramIDRequest{
		TelegramId: clientID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.Unauthenticated {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.ErrorInfo:
					reason := t.GetReason()
					if reason == "user not found" {
						name := update.Message.Chat.FirstName
						if name == "" {
							name = update.Message.Chat.UserName
						}

						_, err = b.registratorServiceClient.CreateUser(ctx, &registrator.CreateUserRequest{
							Name:       name,
							TelegramId: clientID,
							State:      registrator.UserState_USER_STATE_WELCOME,
						})
						if err != nil {
							logger.Errorf(ctx, "error while create user: %s", err.Error())
						}

						welcomeMessage := welcomeMessage(ctx, clientID, name)
						_, err = b.bot.Send(welcomeMessage)
						return err
					}
				}
			}
		}

		return err
	}

	switch resp.GetUser().GetState() {
	case registrator.UserState_USER_STATE_CHANGING_EMAIL:
		_, err = b.registratorServiceClient.UpdateUserEmail(ctx, &registrator.UpdateUserEmailRequest{
			UserId: resp.GetUser().GetId(),
			Email:  update.Message.Text,
		})
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, getTranslator(emailChangedLexeme)(ctx))
		_, err = b.bot.Send(msg)
	case registrator.UserState_USER_STATE_CHANGINE_NAME:
		_, err = b.registratorServiceClient.UpdateUserName(ctx, &registrator.UpdateUserNameRequest{
			UserId: resp.GetUser().GetId(),
			Name:   update.Message.Text,
		})
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, getTranslator(nameChangedLexeme)(ctx))
		_, err = b.bot.Send(msg)
	case registrator.UserState_USER_STATE_CHANGING_PHONE:
		_, err = b.registratorServiceClient.UpdateUserPhone(ctx, &registrator.UpdateUserPhoneRequest{
			UserId: resp.GetUser().GetId(),
			Phone:  update.Message.Text,
		})
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, getTranslator(phoneChangedLexeme)(ctx))
		_, err = b.bot.Send(msg)
	default:
		if update.Message.PinnedMessage != nil {
			return nil
		}
	}

	return err
}

func (b *Bot) getGamesWithPhotosMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID
	gamesWithPhotosListLimit := uint32(config.GetValue("GamesWithPhotosListLimit").Uint64())

	resp, err := b.photographerServiceClient.GetGamesWithPhotos(ctx, &registrator.GetGamesWithPhotosRequest{
		Limit:  gamesWithPhotosListLimit,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	if resp.GetTotal() == 0 {
		return tgbotapi.NewMessage(clientID, getTranslator(listOfGamesWithPhotosIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, pbGame := range resp.GetGames() {
		leagueResp, err := b.registratorServiceClient.GetLeagueByID(ctx, &registrator.GetLeagueByIDRequest{
			Id: pbGame.GetLeagueId(),
		})
		if err != nil {
			return nil, err
		}

		placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
			Id: pbGame.GetPlaceId(),
		})
		if err != nil {
			return nil, err
		}

		pbReq := &registrator.GetPhotosByGameIDRequest{
			GameId: pbGame.GetId(),
		}

		request, err := getRequest(ctx, CommandGetGamePhotos, pbReq)
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)
		text := fmt.Sprintf(gamePhotosInfoFormatString, model.ResultPlace(pbGame.GetResultPlace()).String(), leagueResp.GetLeague().GetShortName(), pbGame.GetNumber(), placeResp.GetPlace().GetShortName(), model.DateTime(pbGame.GetDate().AsTime()))

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	leftNext := uint32(0)
	if resp.GetTotal() > gamesWithPhotosListLimit {
		leftNext = resp.GetTotal() - gamesWithPhotosListLimit
	}

	if leftNext > 0 {
		pbReq := &registrator.GetGamesWithPhotosRequest{
			Limit:  gamesWithPhotosListLimit,
			Offset: gamesWithPhotosListLimit,
		}

		request, err := getRequest(ctx, CommandGetListGamesWithPhotosNextPage, pbReq)
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)

		btnNext := tgbotapi.InlineKeyboardButton{
			Text:         nextPageStringText,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNext))
	}

	msg := tgbotapi.NewMessage(clientID, getTranslator(gamePhotosLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (b *Bot) getListOfGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	games, err := b.gamesFacade.GetGames(ctx, true)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return tgbotapi.NewMessage(clientID, getTranslator(listOfGamesIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		pbReq := &registrator.GetGameByIDRequest{
			GameId: game.ID,
		}

		request, err := getRequest(ctx, CommandGetGame, pbReq)
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)

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

	msg := tgbotapi.NewMessage(clientID, getTranslator(chooseGameLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (b *Bot) getListOfMyGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	respUser, err := b.registratorServiceClient.GetUserByTelegramID(ctx, &registrator.GetUserByTelegramIDRequest{
		TelegramId: clientID,
	})
	if err != nil {
		return nil, err
	}

	resp, err := b.registratorServiceClient.GetUserGames(ctx, &registrator.GetUserGamesRequest{
		Active: true,
		UserId: respUser.GetUser().GetId(),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.GetGames()) == 0 {
		return tgbotapi.NewMessage(clientID, getTranslator(listOfMyGamesIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, pbGame := range resp.GetGames() {
		leagueResp, err := b.registratorServiceClient.GetLeagueByID(ctx, &registrator.GetLeagueByIDRequest{
			Id: pbGame.GetLeagueId(),
		})
		if err != nil {
			return nil, err
		}

		placeResp, err := b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
			Id: pbGame.GetPlaceId(),
		})
		if err != nil {
			return nil, err
		}

		pbReq := &registrator.GetGameByIDRequest{
			GameId: pbGame.GetId(),
		}

		request, err := getRequest(ctx, CommandGetGame, pbReq)
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)
		text := fmt.Sprintf(gameInfoFormatString, leagueResp.GetLeague().GetShortName(), pbGame.GetNumber(), placeResp.GetPlace().GetShortName(), model.DateTime(pbGame.GetDate().AsTime()))

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg := tgbotapi.NewMessage(clientID, getTranslator(listOfMyGamesLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (b *Bot) getListOfRegisteredGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	games, err := b.gamesFacade.GetRegisteredGames(ctx, true)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return tgbotapi.NewMessage(clientID, getTranslator(listOfRegisteredGamesIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		pbReq := &registrator.GetGameByIDRequest{
			GameId: game.ID,
		}

		request, err := getRequest(ctx, CommandGetGame, pbReq)
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)
		text := fmt.Sprintf(gameInfoFormatString, game.League.ShortName, game.Number, game.Place.ShortName, game.DateTime())

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg := tgbotapi.NewMessage(clientID, getTranslator(listOfRegisteredGamesLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (b *Bot) getSettingsMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	user, err := b.registratorServiceClient.GetUserByTelegramID(ctx, &registrator.GetUserByTelegramIDRequest{
		TelegramId: clientID,
	})
	if err != nil {
		return nil, err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	{
		request, err := getRequest(ctx, CommandChangeEmail, "")
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)
		btnEmail := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, getTranslator(changeEmailLexeme)(ctx), user.GetUser().GetEmail()),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnEmail))
	}

	{
		request, err := getRequest(ctx, CommandChangeName, "")
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)
		btnName := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, getTranslator(changeNameLexeme)(ctx), user.GetUser().GetName()),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnName))
	}

	{
		request, err := getRequest(ctx, CommandChangePhone, "")
		if err != nil {
			return nil, err
		}

		callbackData := b.registerRequest(ctx, request)
		btnPhone := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, getTranslator(changePhoneLexeme)(ctx), user.GetUser().GetPhone()),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPhone))
	}

	msg := tgbotapi.NewMessage(clientID, getTranslator(settingsLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func helpMessage(ctx context.Context, clientID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(clientID, getTranslator(helpMessageLexeme)(ctx))

	return msg
}

func welcomeMessage(ctx context.Context, clientID int64, name string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(clientID, fmt.Sprintf(getTranslator(welcomeMessageLexeme)(ctx), name))

	return msg
}
