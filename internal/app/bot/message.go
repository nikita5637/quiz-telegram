package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/users"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	callbackdata_utils "github.com/nikita5637/quiz-telegram/internal/pkg/utils/callbackdata"
	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"
	userutils "github.com/nikita5637/quiz-telegram/utils/user"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/commands"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/icons"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
)

const (
	gameInfoFormatString       = "%s %s %s %s\n"
	gamePhotosInfoFormatString = icons.Photo + "%s" + gameInfoFormatString
	settingFormatString        = "%s [%s]"
)

var (
	birthdateChangedLexeme = i18n.Lexeme{
		Key:      "birthdate_changed",
		FallBack: "Birthdate changed",
	}
	certificateInfoLexeme = i18n.Lexeme{
		Key:      "certificate_info",
		FallBack: "Certificate info",
	}
	certificateTypeLexeme = i18n.Lexeme{
		Key:      "certificate_type",
		FallBack: "Certificate type",
	}
	changeBirthdateLexeme = i18n.Lexeme{
		Key:      "change_birthdate",
		FallBack: "Change birthdate",
	}
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
	changeSexLexeme = i18n.Lexeme{
		Key:      "change_sex",
		FallBack: "Change sex",
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
	listOfCertificatesIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_certificates_is_empty",
		FallBack: "There are no certificates",
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
	myGamesLexeme = i18n.Lexeme{
		Key:      "my_games",
		FallBack: "My games",
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
	registeredGamesLexeme = i18n.Lexeme{
		Key:      "registered_games",
		FallBack: "Registered games",
	}
	settingsLexeme = i18n.Lexeme{
		Key:      "settings",
		FallBack: "Settings",
	}
	sexChangedLexeme = i18n.Lexeme{
		Key:      "sex_changed",
		FallBack: "Sex changed",
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

	user, err := b.checkAuth(ctx, clientID)
	if err != nil {
		name := firstName
		if name == "" {
			name = userName
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

	ctx = userutils.NewContextWithUser(ctx, user)

	var handler func(ctx context.Context) error

	switch text {
	case "/certificates":
		handler = func(ctx context.Context) error {
			var msg tgbotapi.Chattable
			msg, err = b.getListOfCertificatesMessage(ctx, update)
			if err != nil {
				return err
			}

			_, err = b.bot.Send(msg)
			if err != nil {
				logger.Errorf(ctx, "error while sending message: %w", err)
			}

			return err
		}
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
	case "/mygames", i18n.GetTranslator(myGamesLexeme)(ctx):
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
	case "/registeredgames", i18n.GetTranslator(registeredGamesLexeme)(ctx):
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
	case "/settings", i18n.GetTranslator(settingsLexeme)(ctx):
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
		responseMessage := tgbotapi.NewMessage(clientID, i18n.GetTranslator(somethingWentWrongLexeme)(ctx))
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.PermissionDenied {
				for _, detail := range st.Details() {
					switch t := detail.(type) {
					case *errdetails.ErrorInfo:
						reason := t.GetReason()
						if reason == "banned" {
							responseMessage = tgbotapi.NewMessage(clientID, i18n.GetTranslator(permissionDeniedLexeme)(ctx))
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
	user, err := b.usersFacade.GetUserByTelegramID(ctx, clientID)
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.Unauthenticated {
			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.ErrorInfo:
					reason := t.GetReason()
					if reason == users.ReasonUserNotFound {
						name := update.Message.Chat.FirstName
						if name == "" {
							name = update.Message.Chat.UserName
						}

						_, err = b.usersFacade.CreateUser(ctx, name, clientID, int32(usermanagerpb.UserState_USER_STATE_WELCOME))
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

	switch user.State {
	case int32(usermanagerpb.UserState_USER_STATE_CHANGING_BIRTHDATE):
		err = b.usersFacade.UpdateUserBirthdate(ctx, user.ID, update.Message.Text)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(birthdateChangedLexeme)(ctx))
		msg.ReplyMarkup = replyKeyboardMarkup(ctx)

		_, err = b.bot.Send(msg)
	case int32(usermanagerpb.UserState_USER_STATE_CHANGING_EMAIL):
		err = b.usersFacade.UpdateUserEmail(ctx, user.ID, update.Message.Text)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(emailChangedLexeme)(ctx))
		msg.ReplyMarkup = replyKeyboardMarkup(ctx)

		_, err = b.bot.Send(msg)
	case int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME):
		err = b.usersFacade.UpdateUserName(ctx, user.ID, update.Message.Text)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(nameChangedLexeme)(ctx))
		msg.ReplyMarkup = replyKeyboardMarkup(ctx)

		_, err = b.bot.Send(msg)
	case int32(usermanagerpb.UserState_USER_STATE_CHANGING_PHONE):
		err = b.usersFacade.UpdateUserPhone(ctx, user.ID, update.Message.Text)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(phoneChangedLexeme)(ctx))
		msg.ReplyMarkup = replyKeyboardMarkup(ctx)

		_, err = b.bot.Send(msg)
	case int32(usermanagerpb.UserState_USER_STATE_CHANGING_SEX):
		err = b.usersFacade.UpdateUserSex(ctx, user.ID, model.SexFromString(update.Message.Text))
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(sexChangedLexeme)(ctx))
		msg.ReplyMarkup = replyKeyboardMarkup(ctx)

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

	games, total, err := b.gamePhotosFacade.GetGamesWithPhotos(ctx, gamesWithPhotosListLimit, 0)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfGamesWithPhotosIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &commands.GetGamePhotosData{
			GameID: game.ID,
		}

		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandGetGamePhotos, payload)
		if err != nil {
			return nil, err
		}

		text := fmt.Sprintf(gamePhotosInfoFormatString, game.ResultPlace.String(), game.League.ShortName, game.Number, game.Place.ShortName, game.DateTime())

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	leftNext := uint32(0)
	if total > gamesWithPhotosListLimit {
		leftNext = total - gamesWithPhotosListLimit
	}

	if leftNext > 0 {
		payload := &commands.GetGamesWithPhotosData{
			Limit:  gamesWithPhotosListLimit,
			Offset: gamesWithPhotosListLimit,
		}

		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandGetListGamesWithPhotosNextPage, payload)
		if err != nil {
			return nil, err
		}

		btnNext := tgbotapi.InlineKeyboardButton{
			Text:         icons.NextPage,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNext))
	}

	msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(gamePhotosLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (b *Bot) getListOfCertificatesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	certificates, err := b.certificatesFacade.GetActiveCertificates(ctx)
	if err != nil {
		return nil, err
	}

	if len(certificates) == 0 {
		return tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfCertificatesIsEmptyLexeme)(ctx)), nil
	}

	textBuilder := strings.Builder{}
	for _, certificate := range certificates {
		textBuilder.WriteString(fmt.Sprintf("%s: %s\n", i18n.GetTranslator(certificateTypeLexeme)(ctx), certificate.Type))
		textBuilder.WriteString(fmt.Sprintf("%s: %s\n", i18n.GetTranslator(certificateInfoLexeme)(ctx), certificate.Info))
		textBuilder.WriteString("\n")
	}

	msg := tgbotapi.NewMessage(clientID, textBuilder.String())

	return msg, nil
}

func (b *Bot) getListOfGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	games, err := b.gamesFacade.GetGames(ctx, true)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfGamesIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &commands.GetGameData{
			GameID: game.ID,
		}

		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandGetGame, payload)
		if err != nil {
			return nil, err
		}

		text := fmt.Sprintf(gameInfoFormatString, game.League.ShortName, game.Number, game.Place.ShortName, game.Date)

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

	msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(chooseGameLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (b *Bot) getListOfMyGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	user, err := b.usersFacade.GetUserByTelegramID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	games, err := b.gamesFacade.GetUserGames(ctx, true, user.ID)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfMyGamesIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &commands.GetGameData{
			GameID: game.ID,
		}

		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandGetGame, payload)
		if err != nil {
			return nil, err
		}

		text := fmt.Sprintf(gameInfoFormatString, game.League.ShortName, game.Number, game.Place.ShortName, game.DateTime())

		btn := tgbotapi.InlineKeyboardButton{
			Text:         text,
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfMyGamesLexeme)(ctx))
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
		return tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfRegisteredGamesIsEmptyLexeme)(ctx)), nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, game := range games {
		payload := &commands.GetGameData{
			GameID: game.ID,
		}

		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandGetGame, payload)
		if err != nil {
			return nil, err
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

	msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfRegisteredGamesLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (b *Bot) getSettingsMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := update.Message.From.ID

	user, err := b.usersFacade.GetUserByTelegramID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	{
		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandChangeEmail, "")
		if err != nil {
			return nil, err
		}

		btnEmail := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changeEmailLexeme)(ctx), user.Email),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnEmail))
	}

	{
		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandChangeName, "")
		if err != nil {
			return nil, err
		}

		btnName := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changeNameLexeme)(ctx), user.Name),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnName))
	}

	{
		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandChangePhone, "")
		if err != nil {
			return nil, err
		}

		btnPhone := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changePhoneLexeme)(ctx), user.Phone),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPhone))
	}

	{
		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandChangeBirthdate, "")
		if err != nil {
			return nil, err
		}

		btnBirthdate := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changeBirthdateLexeme)(ctx), user.Birthdate),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnBirthdate))
	}

	{
		callbackData, err := callbackdata_utils.GetCallbackData(ctx, commands.CommandChangeSex, "")
		if err != nil {
			return nil, err
		}

		btnSex := tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changeSexLexeme)(ctx), user.Sex),
			CallbackData: &callbackData,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnSex))
	}

	msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(settingsLexeme)(ctx))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func helpMessage(ctx context.Context, clientID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(helpMessageLexeme)(ctx))

	msg.ReplyMarkup = replyKeyboardMarkup(ctx)

	return msg
}

func welcomeMessage(ctx context.Context, clientID int64, name string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(clientID, fmt.Sprintf(i18n.GetTranslator(welcomeMessageLexeme)(ctx), name))

	return msg
}
