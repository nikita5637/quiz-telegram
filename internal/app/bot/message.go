package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	certificatemanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/certificate_manager"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/commands"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/icons"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	callbackdatautils "github.com/nikita5637/quiz-telegram/internal/pkg/utils/callbackdata"
	telegramutils "github.com/nikita5637/quiz-telegram/utils/telegram"
	userutils "github.com/nikita5637/quiz-telegram/utils/user"
	"github.com/spf13/viper"
)

type certificateInfo struct {
	Sum        uint16 `json:"sum,omitempty"`
	Person     uint8  `json:"person,omitempty"`
	ValidFrom  string `json:"valid_from,omitempty"`
	ValidUntil string `json:"valid_until,omitempty"`
}

const (
	gameInfoFormatString         = "%s %s %s %s\n"
	extendedGameInfoFormatString = "%s%s" + gameInfoFormatString
	settingFormatString          = "%s [%s]"
)

var (
	barLexeme = i18n.Lexeme{
		Key:      "bar",
		FallBack: "Bar",
	}
	barBillPaymentLexeme = i18n.Lexeme{
		Key:      "bar_bill_payment",
		FallBack: "Bar bill payment",
	}
	birthdateChangedLexeme = i18n.Lexeme{
		Key:      "birthdate_changed",
		FallBack: "Birthdate changed",
	}
	buyElephantLexeme = i18n.Lexeme{
		Key:      "buy_elephant",
		FallBack: "Everybody says \"%s\", you buy an elephant",
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
	freePassLexeme = i18n.Lexeme{
		Key:      "free_pass",
		FallBack: "Free pass",
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
	listOfPassedGamesIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_passed_games_is_empty",
		FallBack: "There are not passed games",
	}
	listOfRegisteredGamesLexeme = i18n.Lexeme{
		Key:      "list_of_registered_games",
		FallBack: "List of registered games",
	}
	listOfRegisteredGamesIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_registered_games_is_empty",
		FallBack: "There are not registered games",
	}
	listOfYourGamesLexeme = i18n.Lexeme{
		Key:      "list_of_your_games",
		FallBack: "List of your games",
	}
	listOfYourGamesIsEmptyLexeme = i18n.Lexeme{
		Key:      "list_of_your_games_is_empty",
		FallBack: "You don't play with us yet",
	}
	myGamesLexeme = i18n.Lexeme{
		Key:      "my_games",
		FallBack: "My games",
	}
	nameChangedLexeme = i18n.Lexeme{
		Key:      "name_changed",
		FallBack: "Name changed",
	}
	numberOfPersonsLexeme = i18n.Lexeme{
		Key:      "number_of_persons",
		FallBack: "Number of persons",
	}
	passedGamesLexeme = i18n.Lexeme{
		Key:      "passed_games",
		FallBack: "Passed games",
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
	sumLexeme = i18n.Lexeme{
		Key:      "sum",
		FallBack: "Sum",
	}
	validUntilLexeme = i18n.Lexeme{
		Key:      "valid_until",
		FallBack: "Valid until",
	}
)

func (b *Bot) handleMessage(ctx context.Context, update *tgbotapi.Update) error {
	text := update.Message.Text
	user := userutils.GetUserFromContext(ctx)
	logger.DebugKV(ctx, "new private message incoming", "user", user, "text", text)

	var messageHandler func(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error)
	switch text {
	case "/certificates":
		messageHandler = func(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := b.getListOfCertificatesMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("getting list of certificates message error: %w", err)
			}

			return msg, nil
		}
	case "/games":
		messageHandler = func(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := b.getListOfGamesMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("getting list of games message error: %w", err)
			}

			return msg, nil
		}
	case "/help":
		messageHandler = func(ctx context.Context, udpate *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := getHelpMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("getting help message error: %w", err)
			}

			return msg, nil
		}
	case "/mygames", i18n.GetTranslator(myGamesLexeme)(ctx):
		messageHandler = func(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := b.getListOfUserGamesMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("getting list of user games message error: %w", err)
			}

			return msg, nil
		}
	case "/passedgames", i18n.GetTranslator(passedGamesLexeme)(ctx):
		messageHandler = func(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := b.getListOfPassedAndRegisteredGamesMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("getting list of passed and registered games message error: %w", err)
			}

			return msg, nil
		}
	case "/registeredgames", i18n.GetTranslator(registeredGamesLexeme)(ctx):
		messageHandler = func(ctx context.Context, udpate *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := b.getListOfRegisteredGamesMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("getting list of registered games message error: %w", err)
			}

			return msg, nil
		}
	case "/settings", i18n.GetTranslator(settingsLexeme)(ctx):
		messageHandler = func(ctx context.Context, udpate *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := b.getSettingsMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("getting settings message error: %w", err)
			}

			return msg, nil
		}
	default:
		messageHandler = func(ctx context.Context, udpate *tgbotapi.Update) (tgbotapi.Chattable, error) {
			msg, err := b.handleDefaultMessage(ctx, update)
			if err != nil {
				return nil, fmt.Errorf("handling default message error: %w", err)
			}

			return msg, nil
		}
	}

	if messageHandler != nil {
		msg, err := messageHandler(ctx, update)
		if err != nil {
			return fmt.Errorf("messageHandler error: %w", err)
		}

		if messageConfig, ok := msg.(*tgbotapi.MessageConfig); ok {
			if messageConfig.ReplyMarkup == nil {
				kb := tgbotapi.NewReplyKeyboard(
					[]tgbotapi.KeyboardButton{
						tgbotapi.NewKeyboardButton(i18n.GetTranslator(myGamesLexeme)(ctx)),
						tgbotapi.NewKeyboardButton(i18n.GetTranslator(registeredGamesLexeme)(ctx)),
					},
					[]tgbotapi.KeyboardButton{
						tgbotapi.NewKeyboardButton(i18n.GetTranslator(passedGamesLexeme)(ctx)),
					},
					[]tgbotapi.KeyboardButton{
						tgbotapi.NewKeyboardButton(i18n.GetTranslator(settingsLexeme)(ctx)),
					},
				)
				kb.ResizeKeyboard = true

				messageConfig.ReplyMarkup = kb
				logger.Debug(ctx, "added reply keyboard to message")
			}

			if _, err = b.bot.Send(messageConfig); err != nil {
				return fmt.Errorf("sending message error: %w", err)
			}
		} else {
			if _, err = b.bot.Send(msg); err != nil {
				return fmt.Errorf("sending message error: %w", err)
			}
		}
	}

	return nil
}

func (b *Bot) getListOfCertificatesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
		clientID := telegramutils.ClientIDFromContext(ctx)

		certificates, err := b.certificatesFacade.GetActiveCertificates(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting active certificates error: %w", err)
		}

		if len(certificates) == 0 {
			msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfCertificatesIsEmptyLexeme)(ctx))
			return &msg, nil
		}

		textBuilder := strings.Builder{}
		for _, certificate := range certificates {
			var certInfo certificateInfo
			if err := json.Unmarshal([]byte(certificate.Info), &certInfo); err != nil {
				return nil, fmt.Errorf("unmarshaling certificate info error: %w", err)
			}

			wonOnGame, err := b.gamesFacade.GetGame(ctx, certificate.WonOn)
			if err != nil {
				return nil, fmt.Errorf("getting game error: %w", err)
			}

			wonOnGamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, wonOnGame.ID)
			if err != nil {
				return nil, fmt.Errorf("getting game players by game ID error: %w", err)
			}

			place, err := b.placesFacade.GetPlace(ctx, wonOnGame.PlaceID)
			if err != nil {
				return nil, fmt.Errorf("getting place error: %w", err)
			}

			if certificate.Type == int32(certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_FREE_PASS) {
				textBuilder.WriteString(fmt.Sprintf("%s: %s\n", i18n.GetTranslator(certificateTypeLexeme)(ctx), i18n.GetTranslator(freePassLexeme)(ctx)))
				textBuilder.WriteString(fmt.Sprintf("%s: %d\n", i18n.GetTranslator(numberOfPersonsLexeme)(ctx), certInfo.Person))
			} else if certificate.Type == int32(certificatemanagerpb.CertificateType_CERTIFICATE_TYPE_BAR_BILL_PAYMENT) {
				textBuilder.WriteString(fmt.Sprintf("%s: %s\n", i18n.GetTranslator(certificateTypeLexeme)(ctx), i18n.GetTranslator(barBillPaymentLexeme)(ctx)))
				textBuilder.WriteString(fmt.Sprintf("%s: %s\n", i18n.GetTranslator(barLexeme)(ctx), place.Name))
				textBuilder.WriteString(fmt.Sprintf("%s: %d₽\n", i18n.GetTranslator(sumLexeme)(ctx), certInfo.Sum))
			}

			if certInfo.ValidUntil != "" {
				validUntilDate, err := time.Parse("2006-01-02", certInfo.ValidUntil)
				if err != nil {
					logger.Warnf(ctx, "parsing certificate valid until date error: %s", err.Error())
					textBuilder.WriteString(fmt.Sprintf("%s: %s\n", i18n.GetTranslator(validUntilLexeme)(ctx), i18n.GetTranslator(unknownLexeme)(ctx)))
				} else {
					textBuilder.WriteString(fmt.Sprintf("%s: %s\n", i18n.GetTranslator(validUntilLexeme)(ctx), validUntilDate.Format("02.01.2006")))
				}
			}

			playerNames := make([]string, 0)
			for _, wonOnGamePlayer := range wonOnGamePlayers {
				if userID, isPresent := wonOnGamePlayer.UserID.Get(); isPresent {
					user, err := b.usersFacade.GetUser(ctx, userID)
					if err != nil {
						return nil, fmt.Errorf("getting user error: %w", err)
					}

					playerNames = append(playerNames, user.Name)
				}
			}

			if len(playerNames) > 0 {
				textBuilder.WriteString(fmt.Sprintf("%s:\n", i18n.GetTranslator(listOfPlayersLexeme)(ctx)))
				for i, playerName := range playerNames {
					textBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, playerName))
				}
			}

			textBuilder.WriteString("-------------------\n")
		}

		msg := tgbotapi.NewMessage(clientID, textBuilder.String())
		return &msg, nil
	}

	msg, err := fn(ctx, update)
	if err != nil {
		return nil, fmt.Errorf("preparing certificates list message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) getListOfGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
		user := userutils.GetUserFromContext(ctx)

		games, err := b.gamesFacade.GetGames(ctx, false, true, false)
		if err != nil {
			return nil, fmt.Errorf("getting games error: %w", err)
		}

		if len(games) == 0 {
			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(listOfGamesIsEmptyLexeme)(ctx))
			return &msg, nil
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, game := range games {
			league, err := b.leaguesFacade.GetLeague(ctx, game.LeagueID)
			if err != nil {
				return nil, fmt.Errorf("getting league error: %w", err)
			}

			place, err := b.placesFacade.GetPlace(ctx, game.PlaceID)
			if err != nil {
				return nil, fmt.Errorf("getting place error: %w", err)
			}

			text := fmt.Sprintf(gameInfoFormatString, league.ShortName, game.Number, place.ShortName, game.DateTime)

			gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, game.ID)
			if err != nil {
				return nil, fmt.Errorf("getting game players by game ID error: %w", err)
			}

			userWillPlay := false
			numberOfLegioners := 0
			numberOfPlayers := 0
			for _, gamePlayer := range gamePlayers {
				if gamePlayer.UserID.Value() == user.ID {
					userWillPlay = true
				}

				if _, isPresent := gamePlayer.UserID.Get(); isPresent {
					numberOfPlayers++
				} else {
					numberOfLegioners++
				}
			}

			if userWillPlay {
				text = fmt.Sprintf("%s %s", icons.Fist, text)
			} else {
				if numberOfLegioners+numberOfPlayers > 0 {
					text = fmt.Sprintf("%s %s", icons.GameWithPlayers, text)
				}
			}

			payload := &commands.GetGameData{
				GameID:                  game.ID,
				PageIndex:               0,
				GetRootGamesListCommand: commands.CommandGetGamesList,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGame, payload)
			if err != nil {
				return nil, fmt.Errorf("getting callback data error: %w", err)
			}

			btn := tgbotapi.InlineKeyboardButton{
				Text:         text,
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(chooseGameLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		return &msg, nil
	}

	msg, err := fn(ctx, update)
	if err != nil {
		return nil, fmt.Errorf("preparing list of games message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) getListOfPassedAndRegisteredGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
		user := userutils.GetUserFromContext(ctx)

		passedGamesListLimit := viper.GetUint64("bot.passed_games_list_limit")

		passedGames, total, err := b.gamesFacade.SearchPassedAndRegisteredGames(ctx, 1, passedGamesListLimit)
		if err != nil {
			return nil, fmt.Errorf("searching passed and registered games error: %w", err)
		}

		if total == 0 {
			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(listOfPassedGamesIsEmptyLexeme)(ctx))
			return &msg, nil
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, passedGame := range passedGames {
			payload := &commands.GetGameData{
				GameID:    passedGame.ID,
				PageIndex: 0,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGame, payload)
			if err != nil {
				return nil, err
			}

			var resultPlace model.ResultPlace
			gameResult, err := b.gameResultsFacade.GetGameResultByGameID(ctx, passedGame.ID)
			if err != nil {
				logger.ErrorKV(ctx, fmt.Sprintf("getting game result by game ID error: %s", err.Error()), "game", passedGame)
			} else {
				resultPlace = gameResult.ResultPlace
			}

			league, err := b.leaguesFacade.GetLeague(ctx, passedGame.LeagueID)
			if err != nil {
				return nil, fmt.Errorf("getting league error: %w", err)
			}

			place, err := b.placesFacade.GetPlace(ctx, passedGame.PlaceID)
			if err != nil {
				return nil, fmt.Errorf("getting place error: %w", err)
			}

			gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, passedGame.ID)
			if err != nil {
				return nil, fmt.Errorf("getting game players by game ID error: %w", err)
			}

			userHasPlayed := false
			for _, gamePlayer := range gamePlayers {
				if userID, isPresent := gamePlayer.UserID.Get(); isPresent {
					if user.ID == userID {
						userHasPlayed = true
						break
					}
				}
			}

			fist := ""
			if userHasPlayed {
				fist = icons.Fist
			}
			text := fmt.Sprintf(extendedGameInfoFormatString, fist, resultPlace.String(), league.ShortName, passedGame.Number, place.ShortName, passedGame.DateTime)

			btn := tgbotapi.InlineKeyboardButton{
				Text:         text,
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		if total > passedGamesListLimit {
			payload := &commands.GetPassedAndRegisteredGamesListData{
				Page:     2,
				PageSize: passedGamesListLimit,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetPassedAndRegisteredGamesList, payload)
			if err != nil {
				return nil, err
			}

			btnNext := tgbotapi.InlineKeyboardButton{
				Text:         icons.NextPage,
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNext))
		}

		msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(passedGamesLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		return &msg, nil
	}
	msg, err := fn(ctx, update)
	if err != nil {
		return nil, fmt.Errorf("preparing list of passed and registered games message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) getListOfRegisteredGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
		user := userutils.GetUserFromContext(ctx)

		registeredGames, err := b.gamesFacade.GetGames(ctx, true, true, false)
		if err != nil {
			return nil, fmt.Errorf("getting games error: %w", err)
		}

		if len(registeredGames) == 0 {
			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(listOfRegisteredGamesIsEmptyLexeme)(ctx))
			return &msg, nil
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, registeredGame := range registeredGames {
			league, err := b.leaguesFacade.GetLeague(ctx, registeredGame.LeagueID)
			if err != nil {
				return nil, fmt.Errorf("getting league error: %w", err)
			}

			place, err := b.placesFacade.GetPlace(ctx, registeredGame.PlaceID)
			if err != nil {
				return nil, fmt.Errorf("getting place error: %w", err)
			}

			text := fmt.Sprintf(gameInfoFormatString, league.ShortName, registeredGame.Number, place.ShortName, registeredGame.DateTime)

			gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, registeredGame.ID)
			if err != nil {
				return nil, fmt.Errorf("getting game players by game ID error: %w", err)
			}

			userWillPlay := false
			numberOfLegioners := 0
			numberOfPlayers := 0
			for _, gamePlayer := range gamePlayers {
				if gamePlayer.UserID.Value() == user.ID {
					userWillPlay = true
				}

				if _, isPresent := gamePlayer.UserID.Get(); isPresent {
					numberOfPlayers++
				} else {
					numberOfLegioners++
				}
			}

			if userWillPlay {
				text = fmt.Sprintf("%s %s", icons.Fist, text)
			} else {
				if numberOfLegioners+numberOfPlayers > 0 {
					text = fmt.Sprintf("%s %s", icons.GameWithPlayers, text)
				}
			}

			payload := &commands.GetGameData{
				GameID:                  registeredGame.ID,
				PageIndex:               0,
				GetRootGamesListCommand: commands.CommandGetRegisteredGamesList,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGame, payload)
			if err != nil {
				return nil, fmt.Errorf("getting callback data error: %w", err)
			}

			btn := tgbotapi.InlineKeyboardButton{
				Text:         text,
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(listOfRegisteredGamesLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		return &msg, nil
	}

	msg, err := fn(ctx, update)
	if err != nil {
		return nil, fmt.Errorf("preparing list of registered games message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) getListOfUserGamesMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
		user := userutils.GetUserFromContext(ctx)

		userGames, err := b.gamesFacade.GetGamesByUserID(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("getting games by user ID error: %w", err)
		}

		if len(userGames) == 0 {
			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(listOfYourGamesIsEmptyLexeme)(ctx))
			return &msg, nil
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, userGame := range userGames {
			league, err := b.leaguesFacade.GetLeague(ctx, userGame.LeagueID)
			if err != nil {
				return nil, fmt.Errorf("getting league error: %w", err)
			}

			place, err := b.placesFacade.GetPlace(ctx, userGame.PlaceID)
			if err != nil {
				return nil, fmt.Errorf("getting place error: %w", err)
			}

			text := fmt.Sprintf(gameInfoFormatString, league.ShortName, userGame.Number, place.ShortName, userGame.DateTime)
			if !userGame.Registered {
				text = fmt.Sprintf(extendedGameInfoFormatString, "", icons.UnregisteredGame, league.ShortName, userGame.Number, place.ShortName, userGame.DateTime)
			}

			payload := &commands.GetGameData{
				GameID:                  userGame.ID,
				PageIndex:               0,
				GetRootGamesListCommand: commands.CommandGetUserGamesList,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGame, payload)
			if err != nil {
				return nil, fmt.Errorf("getting callback data error: %w", err)
			}

			btn := tgbotapi.InlineKeyboardButton{
				Text:         text,
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(listOfYourGamesLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		return &msg, nil
	}

	msg, err := fn(ctx, update)
	if err != nil {
		return nil, fmt.Errorf("preparing list of user games message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) getSettingsMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
		user := userutils.GetUserFromContext(ctx)

		rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
		{
			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandChangeEmail, "")
			if err != nil {
				return nil, err
			}

			btnEmail := tgbotapi.InlineKeyboardButton{
				Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changeEmailLexeme)(ctx), user.Email.Value()),
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnEmail))
		}

		{
			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandChangeName, "")
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
			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandChangePhone, "")
			if err != nil {
				return nil, err
			}

			btnPhone := tgbotapi.InlineKeyboardButton{
				Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changePhoneLexeme)(ctx), user.Phone.Value()),
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPhone))
		}

		{
			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandChangeBirthdate, "")
			if err != nil {
				return nil, err
			}

			btnBirthdate := tgbotapi.InlineKeyboardButton{
				Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changeBirthdateLexeme)(ctx), user.Birthdate.Value()),
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnBirthdate))
		}

		{
			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandChangeSex, "")
			if err != nil {
				return nil, err
			}

			btnSex := tgbotapi.InlineKeyboardButton{
				Text:         fmt.Sprintf(settingFormatString, i18n.GetTranslator(changeSexLexeme)(ctx), user.Sex.Value()),
				CallbackData: &callbackData,
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnSex))
		}

		msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(settingsLexeme)(ctx))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		return &msg, nil
	}

	msg, err := fn(ctx, update)
	if err != nil {
		return nil, fmt.Errorf("preparing settings message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) handleDefaultMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
		user := userutils.GetUserFromContext(ctx)

		switch user.State {
		case int32(usermanagerpb.UserState_USER_STATE_CHANGING_BIRTHDATE):
			err := b.usersFacade.UpdateUserBirthdate(ctx, user.ID, update.Message.Text)
			if err != nil {
				return nil, fmt.Errorf("updating user birthdate error: %w", err)
			}

			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(birthdateChangedLexeme)(ctx))
			return &msg, nil
		case int32(usermanagerpb.UserState_USER_STATE_CHANGING_EMAIL):
			err := b.usersFacade.UpdateUserEmail(ctx, user.ID, update.Message.Text)
			if err != nil {
				return nil, fmt.Errorf("updating user email error: %w", err)
			}

			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(emailChangedLexeme)(ctx))
			return &msg, nil
		case int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME):
			err := b.usersFacade.UpdateUserName(ctx, user.ID, update.Message.Text)
			if err != nil {
				return nil, fmt.Errorf("updating user name error: %w", err)
			}

			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(nameChangedLexeme)(ctx))
			return &msg, nil
		case int32(usermanagerpb.UserState_USER_STATE_CHANGING_PHONE):
			err := b.usersFacade.UpdateUserPhone(ctx, user.ID, update.Message.Text)
			if err != nil {
				return nil, fmt.Errorf("updating user phone error: %w", err)
			}

			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(phoneChangedLexeme)(ctx))
			return &msg, nil
		case int32(usermanagerpb.UserState_USER_STATE_CHANGING_SEX):
			err := b.usersFacade.UpdateUserSex(ctx, user.ID, model.SexFromString(update.Message.Text))
			if err != nil {
				return nil, fmt.Errorf("updating user sex error: %w", err)
			}

			msg := tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(sexChangedLexeme)(ctx))
			return &msg, nil
		}

		msg := tgbotapi.NewMessage(user.TelegramID, fmt.Sprintf(i18n.GetTranslator(buyElephantLexeme)(ctx), update.Message.Text))
		return &msg, nil
	}

	msg, err := fn(ctx, update)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func getHelpMessage(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	clientID := telegramutils.ClientIDFromContext(ctx)
	msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(helpMessageLexeme)(ctx))

	return &msg, nil
}
