package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mono83/maybe"
	croupierpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/croupier"
	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/commands"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gameplayers"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/icons"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	callbackdatautils "github.com/nikita5637/quiz-telegram/internal/pkg/utils/callbackdata"
	telegramutils "github.com/nikita5637/quiz-telegram/utils/telegram"
	userutils "github.com/nikita5637/quiz-telegram/utils/user"
	"github.com/spf13/viper"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

var (
	addressLexeme = i18n.Lexeme{
		Key:      "address",
		FallBack: "Address",
	}
	addToCalendarLexeme = i18n.Lexeme{
		Key:      "add_to_calendar",
		FallBack: "Add to calendar",
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
	leagueLexeme = i18n.Lexeme{
		Key:      "league",
		FallBack: "League",
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
	menuLexeme = i18n.Lexeme{
		Key:      "menu",
		FallBack: "Menu",
	}
	mixLexeme = i18n.Lexeme{
		Key:      "mix",
		FallBack: "Mix",
	}
	noFreeSlotLexeme = i18n.Lexeme{
		Key:      "no_free_slot",
		FallBack: "There are not free slot",
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
	placeLexeme = i18n.Lexeme{
		Key:      "place",
		FallBack: "Place",
	}
	playsUnlikelyLexeme = i18n.Lexeme{
		Key:      "plays_unlikely",
		FallBack: "plays unlikely",
	}
	resultPlaceLexeme = i18n.Lexeme{
		Key:      "result_place",
		FallBack: "Result place",
	}
	roundPointsLexeme = i18n.Lexeme{
		Key:      "round_points",
		FallBack: "Round points",
	}
	thereAreNoRegistrationForTheGameLexeme = i18n.Lexeme{
		Key:      "there_are_no_registration_for_the_game",
		FallBack: "There are no registration for the game",
	}
	thereAreNoYourLegionersRegisteredForTheGameLexeme = i18n.Lexeme{
		Key:      "there_are_no_your_legioners_registered_for_the_game",
		FallBack: "There are no your legioners registered for the game",
	}
	titleLexeme = i18n.Lexeme{
		Key:      "title",
		FallBack: "Title",
	}
	youAreNotRegisteredForTheGameLexeme = i18n.Lexeme{
		Key:      "you_are_not_registered_for_the_game",
		FallBack: "You are not registered for the game",
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

func (b *Bot) handleCallbackQuery(ctx context.Context, update *tgbotapi.Update) error {
	callbackData := update.CallbackData()
	user := userutils.GetUserFromContext(ctx)
	logger.DebugKV(ctx, "new callback query incoming", "user", user, "callbackData", callbackData)

	telegramRequest := commands.TelegramRequest{}
	err := json.Unmarshal([]byte(callbackData), &telegramRequest)
	if err != nil {
		return fmt.Errorf("unmarshaling telegram request error: %w", err)
	}

	var callbackHandler func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error)
	switch telegramRequest.Command {
	case commands.CommandChangeBirthdate:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			return b.handleChangeBirthdate(ctx, update, telegramRequest)
		}
	case commands.CommandChangeEmail:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			return b.handleChangeEmail(ctx, update, telegramRequest)
		}
	case commands.CommandChangeName:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			return b.handleChangeName(ctx, update, telegramRequest)
		}
	case commands.CommandChangePhone:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			return b.handleChangePhone(ctx, update, telegramRequest)
		}
	case commands.CommandChangeSex:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			return b.handleChangeSex(ctx, update, telegramRequest)
		}
	case commands.CommandGetGame:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.GetGameData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleGetGame(ctx, update, data)
		}
	case commands.CommandGetGamePhotos:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.GetGamePhotosData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleGetGamePhotos(ctx, update, data)
		}
	case commands.CommandGetGamesList:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.GetGamesListData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleGetGamesList(ctx, update, data)
		}
	case commands.CommandGetPassedAndRegisteredGamesList:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.GetPassedAndRegisteredGamesListData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleGetPassedAndRegisteredGamesList(ctx, update, data)
		}
	case commands.CommandGetVenue:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.GetVenueData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleGetVenue(ctx, update, data)
		}
	case commands.CommandLottery:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.LotteryData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleLottery(ctx, update, data)
		}
	case commands.CommandPlayersListByGame:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.PlayersListByGameData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handlePlayersList(ctx, update, data)
		}
	case commands.CommandRegisterGame:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.RegisterGameData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleRegisterGame(ctx, update, data)
		}
	case commands.CommandRegisterPlayer:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.RegisterPlayerData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleRegisterPlayer(ctx, update, data)
		}
	case commands.CommandUnregisterGame:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.UnregisterGameData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleUnregisterGame(ctx, update, data)
		}
	case commands.CommandUnregisterPlayer:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.UnregisterPlayerData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleUnregisterPlayer(ctx, update, data)
		}
	case commands.CommandUpdatePayment:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.UpdatePaymentData{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleUpdatePayment(ctx, update, data)
		}
	case commands.CommandUpdatePlayerRegistration:
		callbackHandler = func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
			data := &commands.UpdatePlayerRegistration{}
			if err = json.Unmarshal(telegramRequest.Body, data); err != nil {
				return nil, nil, fmt.Errorf("unmarshaling telegram request body error: %w", err)
			}

			return b.handleUpdatePlayerRegistration(ctx, update, data)
		}
	}

	if callbackHandler != nil {
		messages, callbacks, err := callbackHandler(ctx, update)
		if err != nil {
			return fmt.Errorf("callbackHandler error: %w", err)
		}

		for _, message := range messages {
			if _, err := b.bot.Send(message); err != nil {
				return fmt.Errorf("sending message error: %w", err)
			}
		}

		for _, callback := range callbacks {
			if _, err := b.bot.Request(callback); err != nil {
				return fmt.Errorf("sending callback error: %w", err)
			}
		}

		return nil
	}

	return fmt.Errorf("callback handler not found for command: %d", telegramRequest.Command)
}

func (b *Bot) getGameMenu(ctx context.Context, game model.Game, clientID int64, messageID int, page uint32, rootGamesListCommand commands.Command) (*tgbotapi.EditMessageTextConfig, error) {
	if game.HasPassed {
		return b.getPassedAndRegisteredGameMenuEditMessage(ctx, game, clientID, messageID, rootGamesListCommand)
	}

	switch page {
	case 0:
		return b.getGameMenuFirstPageEditMessage(ctx, game, clientID, messageID, rootGamesListCommand)
	case 1:
		return b.getGameMenuSecondPageEditMessage(ctx, game, clientID, messageID, rootGamesListCommand)
	}

	return nil, nil
}

func (b *Bot) getGameMenuFirstPageEditMessage(ctx context.Context, game model.Game, clientID int64, messageID int, rootGamesListCommand commands.Command) (*tgbotapi.EditMessageTextConfig, error) {
	fn := func(ctx context.Context, game model.Game, clientID int64, messageID int, rootGamesListCommand commands.Command) (*tgbotapi.EditMessageTextConfig, error) {
		rows := make([][]tgbotapi.InlineKeyboardButton, 0)

		lotteryResp, err := b.croupierServiceClient.GetLotteryStatus(ctx, &croupierpb.GetLotteryStatusRequest{
			GameId: game.ID,
		})
		if err != nil {
			logger.Warnf(ctx, "getting lottery status error: %w", err)
		}

		if lotteryResp.GetActive() {
			var btnLottery tgbotapi.InlineKeyboardButton
			if btnLottery, err = b.lotteryButton(ctx, game.ID, game.LeagueID, rootGamesListCommand); err != nil {
				return nil, fmt.Errorf("generating lottery button error: %w", err)
			}

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnLottery))
		}

		user := userutils.GetUserFromContext(ctx)

		place, err := b.placesFacade.GetPlace(ctx, game.PlaceID)
		if err != nil {
			return nil, fmt.Errorf("getting place error: %w", err)
		}

		gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, game.ID)
		if err != nil {
			return nil, fmt.Errorf("getting game players by game ID: %w", err)
		}

		userWillPlay := false
		playerDegree := model.DegreeInvalid
		numberOfLegioners := uint32(0)
		numberOfPlayers := uint32(0)
		numberOfUserLegioners := uint32(0)
		for _, gamePlayer := range gamePlayers {
			if userID, isPresent := gamePlayer.UserID.Get(); isPresent {
				if user.ID == userID {
					playerDegree = gamePlayer.Degree
					userWillPlay = true
				}
				numberOfPlayers++
			} else {
				numberOfLegioners++
				if gamePlayer.RegisteredBy == user.ID {
					numberOfUserLegioners++
				}
			}
		}

		if userWillPlay {
			var btn1 tgbotapi.InlineKeyboardButton
			btn1, err = b.unregisterPlayerButton(ctx, game.ID, user.ID, user.ID, rootGamesListCommand)
			if err != nil {
				return nil, fmt.Errorf("generating unregister player button error: %w", err)
			}

			newDegree := model.DegreeInvalid
			if playerDegree == model.DegreeLikely {
				newDegree = model.DegreeUnlikely
			} else if playerDegree == model.DegreeUnlikely {
				newDegree = model.DegreeLikely
			}

			var btn2 tgbotapi.InlineKeyboardButton
			if btn2, err = b.updatePlayerRegistionButton(ctx, game.ID, user.ID, user.ID, newDegree, rootGamesListCommand); err != nil {
				return nil, fmt.Errorf("generating update player registration button error: %w", err)
			}

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1, btn2))
		}

		if numberOfLegioners+numberOfPlayers < game.MaxPlayers {
			if !userWillPlay {
				var btn1 tgbotapi.InlineKeyboardButton
				btn1, err = b.registerPlayerButton(ctx, game.ID, user.ID, user.ID, model.DegreeLikely, rootGamesListCommand)
				if err != nil {
					return nil, fmt.Errorf("generating register player button error: %w", err)
				}

				var btn2 tgbotapi.InlineKeyboardButton
				btn2, err = b.registerPlayerButton(ctx, game.ID, user.ID, user.ID, model.DegreeUnlikely, rootGamesListCommand)
				if err != nil {
					return nil, fmt.Errorf("generating register player button error: %w", err)
				}
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn1, btn2))
			}

			var btn3 tgbotapi.InlineKeyboardButton
			btn3, err = b.registerPlayerButton(ctx, game.ID, 0, user.ID, model.DegreeLikely, rootGamesListCommand)
			if err != nil {
				return nil, fmt.Errorf("generating register player button error: %w", err)
			}

			var btn4 tgbotapi.InlineKeyboardButton
			btn4, err = b.registerPlayerButton(ctx, game.ID, 0, user.ID, model.DegreeUnlikely, rootGamesListCommand)
			if err != nil {
				return nil, fmt.Errorf("generating register player button error: %w", err)
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn3, btn4))
		}

		if numberOfUserLegioners > 0 {
			var btn tgbotapi.InlineKeyboardButton
			btn, err = b.unregisterPlayerButton(ctx, game.ID, 0, user.ID, rootGamesListCommand)
			if err != nil {
				return nil, fmt.Errorf("generating unregister player button error: %w", err)
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		if numberOfLegioners+numberOfPlayers > 0 {
			var btnPlayersList tgbotapi.InlineKeyboardButton
			btnPlayersList, err = b.playersListButton(ctx, game.ID)
			if err != nil {
				return nil, fmt.Errorf("generating players list button error: %w", err)
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPlayersList))
		}

		if !game.Registered {
			var btnRegisterGame tgbotapi.InlineKeyboardButton
			btnRegisterGame, err = b.registerGameButton(ctx, game.ID, rootGamesListCommand)
			if err != nil {
				return nil, fmt.Errorf("generating register game button error: %w", err)
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnRegisterGame))
		} else {
			var btnUnregisterGame tgbotapi.InlineKeyboardButton
			btnUnregisterGame, err = b.unregisterGameButton(ctx, game.ID, rootGamesListCommand)
			if err != nil {
				return nil, fmt.Errorf("generating unregister game button error: %w", err)
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnUnregisterGame))
		}

		// root games list
		var btnPrevPage tgbotapi.InlineKeyboardButton
		{
			getGamesListData := &commands.GetGamesListData{
				Command: rootGamesListCommand,
			}
			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGamesList, getGamesListData)
			if err != nil {
				return nil, fmt.Errorf("getting callback data error: %w", err)
			}

			btnPrevPage = tgbotapi.InlineKeyboardButton{
				Text:         icons.PrevPage,
				CallbackData: &callbackData,
			}
		}

		// second game menu page
		var btnNextMenuPage tgbotapi.InlineKeyboardButton
		{
			getGameData := &commands.GetGameData{
				GameID:                  game.ID,
				PageIndex:               1,
				GetRootGamesListCommand: rootGamesListCommand,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGame, getGameData)
			if err != nil {
				return nil, fmt.Errorf("getting callback data error: %w", err)
			}

			btnNextMenuPage = tgbotapi.InlineKeyboardButton{
				Text:         icons.NextPage,
				CallbackData: &callbackData,
			}
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPrevPage, btnNextMenuPage))

		msg := tgbotapi.NewEditMessageText(clientID, messageID, getGameMenuInfo(ctx, game, place, numberOfPlayers, numberOfLegioners))
		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &replyMarkup

		return &msg, nil
	}

	msg, err := fn(ctx, game, clientID, messageID, rootGamesListCommand)
	if err != nil {
		return nil, fmt.Errorf("preparing game menu first page edit message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) getGameMenuSecondPageEditMessage(ctx context.Context, game model.Game, clientID int64, messageID int, rootGamesListCommand commands.Command) (*tgbotapi.EditMessageTextConfig, error) {
	fn := func(ctx context.Context, game model.Game, clientID int64, messageID int, rootGamesListCommand commands.Command) (*tgbotapi.EditMessageTextConfig, error) {
		rows := make([][]tgbotapi.InlineKeyboardButton, 0)
		if game.Registered {
			if payment, isPresent := game.Payment.Get(); isPresent {
				var btnNextPayment tgbotapi.InlineKeyboardButton
				btnNextPayment, err := b.nextPaymentButton(ctx, game.ID, payment, rootGamesListCommand)
				if err != nil {
					return nil, fmt.Errorf("generating next payment button error: %w", err)
				}
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnNextPayment))
			}
		}

		gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, game.ID)
		if err != nil {
			return nil, fmt.Errorf("getting game players by game ID: %w", err)
		}

		numberOfLegioners := uint32(0)
		numberOfPlayers := uint32(0)
		for _, gamePlayer := range gamePlayers {
			if _, isPresent := gamePlayer.UserID.Get(); isPresent {
				numberOfPlayers++
			} else {
				numberOfLegioners++
			}
		}

		place, err := b.placesFacade.GetPlace(ctx, game.PlaceID)
		if err != nil {
			return nil, fmt.Errorf("getting place error: %w", err)
		}

		barButtonsRow := []tgbotapi.InlineKeyboardButton{}
		if place.Latitude != 0 && place.Longitude != 0 {
			var btnVenue tgbotapi.InlineKeyboardButton
			btnVenue, err = b.venueButton(ctx, place.ID)
			if err != nil {
				return nil, fmt.Errorf("generating venue button error: %w", err)
			}
			barButtonsRow = append(barButtonsRow, btnVenue)
		}

		if place.MenuLink != "" {
			btnMenu := tgbotapi.NewInlineKeyboardButtonURL(fmt.Sprintf("%s %s", icons.MenuIcon, i18n.GetTranslator(menuLexeme)(ctx)), place.MenuLink)
			barButtonsRow = append(barButtonsRow, btnMenu)
		}

		if len(barButtonsRow) > 0 {
			rows = append(rows, barButtonsRow)
		}

		if game.Registered {
			var icsFile model.ICSFile
			if icsFile, err = b.icsFilesFacade.GetICSFileByGameID(ctx, game.ID); err == nil {
				icsFileButtonsRow := []tgbotapi.InlineKeyboardButton{}
				btn := tgbotapi.NewInlineKeyboardButtonURL(
					i18n.GetTranslator(addToCalendarLexeme)(ctx),
					"http://ics.home0705.keenetic.pro/"+icsFile.Name,
				)
				icsFileButtonsRow = append(icsFileButtonsRow, btn)

				rows = append(rows, icsFileButtonsRow)
			} else {
				logger.Errorf(ctx, "getting ICS file by game ID error: %s", err.Error())
			}
		}

		getGameData := &commands.GetGameData{
			GameID:                  game.ID,
			PageIndex:               0,
			GetRootGamesListCommand: rootGamesListCommand,
		}

		callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGame, getGameData)
		if err != nil {
			return nil, fmt.Errorf("getting callback data error: %w", err)
		}

		btnPrevMenuPage := tgbotapi.InlineKeyboardButton{
			Text:         icons.PrevPage,
			CallbackData: &callbackData,
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPrevMenuPage))

		msg := tgbotapi.NewEditMessageText(clientID, messageID, getGameMenuInfo(ctx, game, place, numberOfPlayers, numberOfLegioners))
		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &replyMarkup

		return &msg, nil
	}

	msg, err := fn(ctx, game, clientID, messageID, rootGamesListCommand)
	if err != nil {
		return nil, fmt.Errorf("preparing game menu second page edit message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) getPassedAndRegisteredGameMenuEditMessage(ctx context.Context, game model.Game, clientID int64, messageID int, rootGamesListCommand commands.Command) (*tgbotapi.EditMessageTextConfig, error) {
	fn := func(ctx context.Context, game model.Game, clientID int64, messageID int, rootGamesListCommand commands.Command) (*tgbotapi.EditMessageTextConfig, error) {
		gameInfo := strings.Builder{}

		league, err := b.leaguesFacade.GetLeague(ctx, game.LeagueID)
		if err != nil {
			return nil, fmt.Errorf("getting league error: %w", err)
		}

		place, err := b.placesFacade.GetPlace(ctx, game.PlaceID)
		if err != nil {
			return nil, fmt.Errorf("getting place error: %w", err)
		}

		gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, game.ID)
		if err != nil {
			return nil, fmt.Errorf("getting game players error: %w", err)
		}

		numberOfPlayers := 0
		numberOfLegioners := 0
		for _, gamePlayer := range gamePlayers {
			if _, isPresent := gamePlayer.UserID.Get(); isPresent {
				numberOfPlayers++
			} else {
				numberOfLegioners++
			}
		}

		gameInfo.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Brain, i18n.GetTranslator(leagueLexeme)(ctx), league.Name))

		if gameName, isPresent := game.Name.Get(); isPresent {
			gameInfo.WriteString(fmt.Sprintf("%s %s: %s %s\n", icons.Sharp, i18n.GetTranslator(titleLexeme)(ctx), gameName, game.Number))
		} else {
			gameInfo.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Sharp, i18n.GetTranslator(numberLexeme)(ctx), game.Number))
		}

		gameInfo.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Calendar, i18n.GetTranslator(dateTimeLexeme)(ctx), game.DateTime))
		gameInfo.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Place, i18n.GetTranslator(placeLexeme)(ctx), place.Name))
		gameInfo.WriteString(fmt.Sprintf("%s %s: %d/%d/%d\n", icons.NumberOfPlayers, i18n.GetTranslator(numberOfPlayersLexeme)(ctx), numberOfPlayers, numberOfLegioners, game.MaxPlayers))

		var gameResult model.GameResult
		if gameResult, err = b.gameResultsFacade.GetGameResultByGameID(ctx, game.ID); err != nil {
			logger.ErrorKV(ctx, fmt.Sprintf("getting game result by game ID error: %s", err.Error()), "game", game)
		} else {
			gameInfo.WriteString(fmt.Sprintf("%s %s: %d\n", icons.StoneFace, i18n.GetTranslator(resultPlaceLexeme)(ctx), gameResult.ResultPlace))
			if roundPointsJSON, isPresent := gameResult.RoundPoints.Get(); isPresent {
				roundPointsMap := map[string]float64{}
				if err = json.Unmarshal([]byte(roundPointsJSON), &roundPointsMap); err != nil {
					logger.ErrorKV(ctx, fmt.Errorf("unmarshaling round points error: %w", err).Error(), "roundPoints", roundPointsJSON)
				} else {
					roundNames := make([]string, 0, len(roundPointsMap))
					for roundName := range roundPointsMap {
						roundNames = append(roundNames, roundName)
					}
					sort.Strings(roundNames)

					gameInfo.WriteString(fmt.Sprintf("%s %s:\n", icons.Info, i18n.GetTranslator(roundPointsLexeme)(ctx)))
					for _, roundName := range roundNames {
						points := roundPointsMap[roundName]
						if float64(int(points)) == points {
							gameInfo.WriteString(fmt.Sprintf("\t%s: %d\n", roundName, int(points)))
						} else {
							gameInfo.WriteString(fmt.Sprintf("\t%s: %.1f\n", roundName, points))
						}
					}
				}
			}
		}

		msg := tgbotapi.NewEditMessageText(clientID, messageID, gameInfo.String())

		rows := make([][]tgbotapi.InlineKeyboardButton, 0)

		btnPlayersList, err := b.playersListButton(ctx, game.ID)
		if err != nil {
			return nil, fmt.Errorf("generating players list button error: %w", err)
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnPlayersList))

		if gamePhotos, err := b.gamePhotosFacade.GetPhotosByGameID(ctx, game.ID); err != nil {
			logger.Errorf(ctx, "getting photos by game ID error: %s", err.Error())
		} else {
			if len(gamePhotos) > 0 {
				btnGamePhotos, err := b.gamePhotosButton(ctx, game.ID)
				if err != nil {
					return nil, fmt.Errorf("generating game photos button error: %w", err)
				}
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnGamePhotos))
			}
		}

		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &replyMarkup

		return &msg, nil
	}

	msg, err := fn(ctx, game, clientID, messageID, rootGamesListCommand)
	if err != nil {
		return nil, fmt.Errorf("preparing passed and registered game menu edit message error: %w", err)
	}

	return msg, nil
}

func (b *Bot) handleChangeBirthdate(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_BIRTHDATE))
}

func (b *Bot) handleChangeEmail(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_EMAIL))
}

func (b *Bot) handleChangeName(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_NAME))
}

func (b *Bot) handleChangePhone(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_PHONE))
}

func (b *Bot) handleChangeSex(ctx context.Context, update *tgbotapi.Update, telegramRequest commands.TelegramRequest) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	return b.updateUserState(ctx, update, int32(usermanagerpb.UserState_USER_STATE_CHANGING_SEX))
}

func (b *Bot) handleGetGame(ctx context.Context, update *tgbotapi.Update, data *commands.GetGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.GetGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		clientID := telegramutils.ClientIDFromContext(ctx)
		messageID := update.CallbackQuery.Message.MessageID
		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		game, err := b.gamesFacade.GetGame(ctx, data.GameID)
		if err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		msg, err := b.getGameMenu(ctx, game, clientID, messageID, data.PageIndex, data.GetRootGamesListCommand)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game menu error: %w", err)
		}

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing game message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleGetGamePhotos(ctx context.Context, update *tgbotapi.Update, data *commands.GetGamePhotosData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.GetGamePhotosData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID
		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		urls, err := b.gamePhotosFacade.GetPhotosByGameID(ctx, data.GameID)
		if err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		messages := []tgbotapi.Chattable{}
		for _, url := range urls {
			messages = append(messages, tgbotapi.NewMessage(clientID, url))
		}

		return messages, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing game photos messages and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleGetGamesList(ctx context.Context, update *tgbotapi.Update, data *commands.GetGamesListData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.GetGamesListData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		callbackQueryID := update.CallbackQuery.ID
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID
		cb := tgbotapi.NewCallback(callbackQueryID, "")

		var msg tgbotapi.Chattable
		switch data.Command {
		case commands.CommandGetRegisteredGamesList:
			var err error
			msg, err = b.getListOfRegisteredGamesMessage(ctx, update)
			if err != nil {
				return nil, nil, fmt.Errorf("getting list of registered games message error: %w", err)
			}
		case commands.CommandGetUserGamesList:
			var err error
			msg, err = b.getListOfUserGamesMessage(ctx, update)
			if err != nil {
				return nil, nil, fmt.Errorf("getting list of user games message error: %w", err)
			}
		default:
			var err error
			msg, err = b.getListOfGamesMessage(ctx, update)
			if err != nil {
				return nil, nil, fmt.Errorf("getting list of games message error: %w", err)
			}
		}

		if m, ok := msg.(*tgbotapi.MessageConfig); ok {
			editMessage := tgbotapi.NewEditMessageText(clientID, messageID, m.Text)
			if replyMarkup, ok := m.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
				editMessage.ReplyMarkup = &replyMarkup
			}
			msg = editMessage
		}

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing games list message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleGetPassedAndRegisteredGamesList(ctx context.Context, update *tgbotapi.Update, data *commands.GetPassedAndRegisteredGamesListData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.GetPassedAndRegisteredGamesListData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		user := userutils.GetUserFromContext(ctx)
		messageID := update.CallbackQuery.Message.MessageID

		passedGamesListLimit := viper.GetUint64("bot.passed_games_list_limit")

		passedGames, total, err := b.gamesFacade.SearchPassedAndRegisteredGames(ctx, data.Page, data.PageSize)
		if err != nil {
			return nil, nil, fmt.Errorf("searching passed and registered games error: %w", err)
		}

		rows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, passedGame := range passedGames {
			payload := &commands.GetGameData{
				GameID:    passedGame.ID,
				PageIndex: 0,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetGame, payload)
			if err != nil {
				return nil, nil, fmt.Errorf("getting callback data error: %w", err)
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
				return nil, nil, fmt.Errorf("getting league error: %w", err)
			}

			place, err := b.placesFacade.GetPlace(ctx, passedGame.PlaceID)
			if err != nil {
				return nil, nil, fmt.Errorf("getting place error: %w", err)
			}

			gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, passedGame.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("getting game players by game ID error: %w", err)
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

		navigateButtonsRow := make([]tgbotapi.InlineKeyboardButton, 0, 2)

		if data.Page > 1 {
			payload := &commands.GetPassedAndRegisteredGamesListData{
				Page:     data.Page - 1,
				PageSize: passedGamesListLimit,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetPassedAndRegisteredGamesList, payload)
			if err != nil {
				return nil, nil, fmt.Errorf("getting callback data error: %w", err)
			}

			btnPrev := tgbotapi.InlineKeyboardButton{
				Text:         icons.PrevPage,
				CallbackData: &callbackData,
			}
			navigateButtonsRow = append(navigateButtonsRow, btnPrev)
		}

		if total > (data.Page * data.PageSize) {
			payload := &commands.GetPassedAndRegisteredGamesListData{
				Page:     data.Page + 1,
				PageSize: passedGamesListLimit,
			}

			callbackData, err := callbackdatautils.GetCallbackData(ctx, commands.CommandGetPassedAndRegisteredGamesList, payload)
			if err != nil {
				return nil, nil, fmt.Errorf("getting callback data error: %w", err)
			}

			btnNext := tgbotapi.InlineKeyboardButton{
				Text:         icons.NextPage,
				CallbackData: &callbackData,
			}
			navigateButtonsRow = append(navigateButtonsRow, btnNext)
		}

		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, navigateButtonsRow)

		msg := tgbotapi.NewEditMessageReplyMarkup(user.TelegramID, messageID, inlineKeyboardMarkup)
		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing passed and registered games list message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleGetVenue(ctx context.Context, update *tgbotapi.Update, data *commands.GetVenueData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.GetVenueData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		clientID := update.CallbackQuery.From.ID

		place, err := b.placesFacade.GetPlace(ctx, data.PlaceID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting place error: %w", err)
		}

		venueConfig := tgbotapi.NewVenue(clientID, place.Name, place.Address, float64(place.Latitude), float64(place.Longitude))
		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		return []tgbotapi.Chattable{venueConfig}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing venue message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleLottery(ctx context.Context, update *tgbotapi.Update, data *commands.LotteryData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.LotteryData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID
		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		resp, err := b.croupierServiceClient.RegisterForLottery(ctx, &croupierpb.RegisterForLotteryRequest{
			GameId: data.GameID,
		})
		if err != nil {
			st := status.Convert(err)

			for _, detail := range st.Details() {
				switch t := detail.(type) {
				case *errdetails.LocalizedMessage:
					localizedMessage := t.GetMessage()
					msg := tgbotapi.NewEditMessageText(clientID, messageID, localizedMessage)
					return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
				}
			}

			return nil, nil, fmt.Errorf("registration for lottery error: %w", err)
		}

		msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(youHaveSuccessfullyRegisteredInLotteryLexeme)(ctx))
		cbs := []tgbotapi.Chattable{cb}
		if resp.GetNumber() > 0 {
			msg = tgbotapi.NewMessage(clientID, fmt.Sprintf("%s: %d", i18n.GetTranslator(yourLotteryNumberIsLexeme)(ctx), resp.GetNumber()))

			unpinMessage := tgbotapi.UnpinAllChatMessagesConfig{
				ChatID: clientID,
			}
			cbs = append(cbs, unpinMessage)

			/*
				pinMessage := tgbotapi.PinChatMessageConfig{
					ChatID:    clientID,
					MessageID: msg.MessageID,
				}
				cbs = append(cbs, pinMessage)
			*/
		}

		return []tgbotapi.Chattable{msg}, cbs, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing lottery message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handlePlayersList(ctx context.Context, update *tgbotapi.Update, data *commands.PlayersListByGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.PlayersListByGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		clientID := update.CallbackQuery.From.ID
		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		gamePlayers, err := b.gamePlayersFacade.GetGamePlayersByGameID(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("geting game players by game ID error: %w", err)
		}

		if len(gamePlayers) == 0 {
			msg := tgbotapi.NewMessage(clientID, i18n.GetTranslator(listOfPlayersIsEmptyLexeme)(ctx))
			return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
		}

		textBuilder := strings.Builder{}
		for i, gamePlayer := range gamePlayers {
			playerName := ""
			if userID, isPresent := gamePlayer.UserID.Get(); isPresent {
				var user model.User
				if user, err = b.usersFacade.GetUser(ctx, userID); err != nil {
					return nil, nil, fmt.Errorf("getting user error: %w", err)
				}
				playerName = user.Name
			} else {
				var user model.User
				if user, err = b.usersFacade.GetUser(ctx, gamePlayer.RegisteredBy); err != nil {
					return nil, nil, fmt.Errorf("getting user error: %w", err)
				}
				playerName = fmt.Sprintf("%s %s", i18n.GetTranslator(legionerByLexeme)(ctx), user.Name)
			}

			if gamePlayer.Degree == model.DegreeUnlikely {
				textBuilder.WriteString(fmt.Sprintf("%d. %s (%s)\n", i+1, playerName, i18n.GetTranslator(playsUnlikelyLexeme)(ctx)))
			} else {
				textBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, playerName))
			}
		}

		msg := tgbotapi.NewMessage(clientID, textBuilder.String())
		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing players list message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleRegisterGame(ctx context.Context, update *tgbotapi.Update, data *commands.RegisterGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.RegisterGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		callbackQueryID := update.CallbackQuery.ID
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID

		if err := b.gamesFacade.RegisterGame(ctx, data.GameID); err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				cb := tgbotapi.NewCallback(callbackQueryID, "")
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("registering game error: %w", err)
		}

		game, err := b.gamesFacade.GetGame(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		msg, err := b.getGameMenu(ctx, game, clientID, messageID, 0, data.GetRootGamesListCommand)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game menu error: %w", err)
		}

		cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(registeredGameLexeme)(ctx))

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing register game message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleRegisterPlayer(ctx context.Context, update *tgbotapi.Update, data *commands.RegisterPlayerData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.RegisterPlayerData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		callbackQueryID := update.CallbackQuery.ID
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID

		gamePlayer := model.GamePlayer{
			GameID:       data.GameID,
			UserID:       maybe.Nothing[int32](),
			RegisteredBy: data.RegisteredBy,
			Degree:       data.Degree,
		}

		if data.UserID != 0 {
			gamePlayer.UserID = maybe.Just(data.UserID)
		}

		if err := b.gamePlayersFacade.RegisterPlayer(ctx, gamePlayer); err != nil {
			if errors.Is(err, games.ErrGameHasPassed) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(games.GameHasPassedLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			} else if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, nil, nil
			} else if errors.Is(err, gameplayers.ErrNoFreeSlot) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(noFreeSlotLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			} else if errors.Is(err, gameplayers.ErrGamePlayerAlreadyRegistered) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreAlreadyRegisteredForTheGameLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			} else if errors.Is(err, gameplayers.ErrThereAreNoRegistrationForTheGame) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(thereAreNoRegistrationForTheGameLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("registering player error: %w", err)
		}

		cb := tgbotapi.NewCallback(callbackQueryID, "")
		if data.UserID == 0 && data.RegisteredBy != data.UserID {
			if data.Degree == model.DegreeLikely {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(legionerIsSignedUpForTheGameLexeme)(ctx))
			} else {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(legionerIsSignedUpForTheGameUnlikelyLexeme)(ctx))
			}
		} else if data.UserID != 0 && data.RegisteredBy == data.UserID {
			if data.Degree == model.DegreeLikely {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreSignedUpForTheGameLexeme)(ctx))
			} else {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreSignedUpForTheGameUnlikelyLexeme)(ctx))
			}
		}

		game, err := b.gamesFacade.GetGame(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		msg, err := b.getGameMenu(ctx, game, clientID, messageID, 0, data.GetRootGamesListCommand)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game menu error: %w", err)
		}

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing register player message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleUnregisterGame(ctx context.Context, update *tgbotapi.Update, data *commands.UnregisterGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.UnregisterGameData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		callbackQueryID := update.CallbackQuery.ID
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID

		if err := b.gamesFacade.UnregisterGame(ctx, data.GameID); err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				cb := tgbotapi.NewCallback(callbackQueryID, "")
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("unregistering game error: %w", err)
		}

		game, err := b.gamesFacade.GetGame(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		msg, err := b.getGameMenu(ctx, game, clientID, messageID, 0, data.GetRootGamesListCommand)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game menu error: %w", err)
		}

		cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(unregisteredGameLexeme)(ctx))

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing unregister game message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleUnregisterPlayer(ctx context.Context, update *tgbotapi.Update, data *commands.UnregisterPlayerData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.UnregisterPlayerData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		callbackQueryID := update.CallbackQuery.ID
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID

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
			if errors.Is(err, games.ErrGameHasPassed) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(games.GameHasPassedLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			} else if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, nil, nil
			} else if errors.Is(err, gameplayers.ErrGamePlayerNotFound) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(thereAreNoYourLegionersRegisteredForTheGameLexeme)(ctx))
				if data.UserID != 0 && data.RegisteredBy == data.UserID {
					cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreNotRegisteredForTheGameLexeme)(ctx))
				}
				return nil, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("unregistering player error: %w", err)
		}

		cb := tgbotapi.NewCallback(callbackQueryID, "")
		if data.UserID == 0 && data.RegisteredBy != data.UserID {
			cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(legionerIsUnsignedUpForTheGameLexeme)(ctx))
		} else if data.UserID != 0 && data.RegisteredBy == data.UserID {
			cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreUnsignedUpForTheGameLexeme)(ctx))
		}

		game, err := b.gamesFacade.GetGame(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		msg, err := b.getGameMenu(ctx, game, clientID, messageID, 0, data.GetRootGamesListCommand)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game menu error: %w", err)
		}

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing unregister player message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleUpdatePayment(ctx context.Context, update *tgbotapi.Update, data *commands.UpdatePaymentData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.UpdatePaymentData) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		callbackQueryID := update.CallbackQuery.ID
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID
		cb := tgbotapi.NewCallback(callbackQueryID, "")

		if err := b.gamesFacade.UpdatePayment(ctx, data.GameID, data.Payment); err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("updating payment error: %w", err)
		}

		game, err := b.gamesFacade.GetGame(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		msg, err := b.getGameMenu(ctx, game, clientID, messageID, 1, data.GetRootGamesListCommand)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game menu error: %w", err)
		}

		switch data.Payment {
		case int32(gamepb.Payment_PAYMENT_CASH):
			cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(cashGamePaymentLexeme)(ctx))
		case int32(gamepb.Payment_PAYMENT_CERTIFICATE):
			cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(freeGamePaymentLexeme)(ctx))
		case int32(gamepb.Payment_PAYMENT_MIXED):
			cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(mixGamePaymentLexeme)(ctx))
		}

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing update payment message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) handleUpdatePlayerRegistration(ctx context.Context, update *tgbotapi.Update, data *commands.UpdatePlayerRegistration) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, data *commands.UpdatePlayerRegistration) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		callbackQueryID := update.CallbackQuery.ID
		clientID := update.CallbackQuery.From.ID
		messageID := update.CallbackQuery.Message.MessageID

		gamePlayer := model.GamePlayer{
			GameID:       data.GameID,
			UserID:       maybe.Just(data.UserID),
			RegisteredBy: data.RegisteredBy,
			Degree:       data.Degree,
		}

		if err := b.gamePlayersFacade.UpdatePlayerRegistration(ctx, gamePlayer); err != nil {
			if errors.Is(err, games.ErrGameNotFound) {
				msg := tgbotapi.NewEditMessageText(clientID, messageID, i18n.GetTranslator(games.GameNotFoundLexeme)(ctx))
				return []tgbotapi.Chattable{msg}, nil, nil
			} else if errors.Is(err, games.ErrGameHasPassed) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(games.GameHasPassedLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			} else if errors.Is(err, gameplayers.ErrGamePlayerNotFound) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreNotRegisteredForTheGameLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			} else if errors.Is(err, gameplayers.ErrThereAreNoRegistrationForTheGame) {
				cb := tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(thereAreNoRegistrationForTheGameLexeme)(ctx))
				return nil, []tgbotapi.Chattable{cb}, nil
			}

			return nil, nil, fmt.Errorf("updateing player registration error: %w", err)
		}

		cb := tgbotapi.NewCallback(callbackQueryID, "")
		if data.UserID == 0 && data.RegisteredBy != data.UserID {
			if data.Degree == model.DegreeUnlikely {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(legionerIsSignedUpForTheGameUnlikelyLexeme)(ctx))
			} else {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(legionerIsSignedUpForTheGameLexeme)(ctx))
			}
		} else if data.UserID != 0 && data.RegisteredBy == data.UserID {
			if data.Degree == model.DegreeUnlikely {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreSignedUpForTheGameUnlikelyLexeme)(ctx))
			} else {
				cb = tgbotapi.NewCallback(callbackQueryID, i18n.GetTranslator(youAreSignedUpForTheGameLexeme)(ctx))
			}
		}

		game, err := b.gamesFacade.GetGame(ctx, data.GameID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game error: %w", err)
		}

		msg, err := b.getGameMenu(ctx, game, clientID, messageID, 0, data.GetRootGamesListCommand)
		if err != nil {
			return nil, nil, fmt.Errorf("getting game menu error: %w", err)
		}

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing update player registration message and callback error: %w", err)
	}

	return msg, cb, nil
}

func (b *Bot) updateUserState(ctx context.Context, update *tgbotapi.Update, state int32) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
	fn := func(ctx context.Context, update *tgbotapi.Update, state int32) ([]tgbotapi.Chattable, []tgbotapi.Chattable, error) {
		user := userutils.GetUserFromContext(ctx)

		if err := b.usersFacade.UpdateUserState(ctx, user.ID, state); err != nil {
			return nil, nil, fmt.Errorf("updating user state errror: %w", err)
		}

		msg := tgbotapi.MessageConfig{}
		switch usermanagerpb.UserState(state) {
		case usermanagerpb.UserState_USER_STATE_CHANGING_BIRTHDATE:
			msg = tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(enterYourBirthdateLexeme)(ctx))
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		case usermanagerpb.UserState_USER_STATE_CHANGING_EMAIL:
			msg = tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(enterYourEmailLexeme)(ctx))
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		case usermanagerpb.UserState_USER_STATE_CHANGING_NAME:
			msg = tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(enterYourNameLexeme)(ctx))
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		case usermanagerpb.UserState_USER_STATE_CHANGING_PHONE:
			msg = tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(enterYourPhoneLexeme)(ctx))
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		case usermanagerpb.UserState_USER_STATE_CHANGING_SEX:
			msg = tgbotapi.NewMessage(user.TelegramID, i18n.GetTranslator(enterYourSexLexeme)(ctx))
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		}

		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		return []tgbotapi.Chattable{msg}, []tgbotapi.Chattable{cb}, nil
	}

	msg, cb, err := fn(ctx, update, state)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing update user state message and callback error: %w", err)
	}

	return msg, cb, nil
}

func getGameMenuInfo(ctx context.Context, game model.Game, place model.Place, numberOfPlayers, numberOfLegioners uint32) string {
	info := strings.Builder{}

	registerStatus := fmt.Sprintf("%s %s", icons.UnregisteredGame, i18n.GetTranslator(unregisteredGameLexeme)(ctx))
	if game.Registered {
		registerStatus = fmt.Sprintf("%s %s", icons.RegisteredGame, i18n.GetTranslator(registeredGameLexeme)(ctx))
	}
	info.WriteString(registerStatus + "\n")

	paymentType := ""
	if gamePaymentType, isPresent := game.PaymentType.Get(); isPresent {
		if strings.Index(gamePaymentType, "cash") != -1 {
			paymentType += strings.ToLower(i18n.GetTranslator(cashLexeme)(ctx))
		}
		if strings.Index(gamePaymentType, "card") != -1 {
			if paymentType != "" {
				paymentType += ", "
			}
			paymentType += strings.ToLower(i18n.GetTranslator(cardLexeme)(ctx))
		}
	}

	if paymentType == "" {
		paymentType = i18n.GetTranslator(unknownLexeme)(ctx)
	}

	if payment, isPresent := game.Payment.Get(); isPresent {
		paymentStatus := fmt.Sprintf("%s %s: %s (%s)", icons.MixGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), strings.ToLower(i18n.GetTranslator(mixLexeme)(ctx)), paymentType)
		if payment == int32(gamepb.Payment_PAYMENT_CASH) {
			paymentStatus = fmt.Sprintf("%s %s: %s", icons.CashGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), paymentType)
		} else if payment == int32(gamepb.Payment_PAYMENT_CERTIFICATE) {
			paymentStatus = fmt.Sprintf("%s %s: %s", icons.FreeGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), strings.ToLower(i18n.GetTranslator(certificateLexeme)(ctx)))
		}
		info.WriteString(paymentStatus + "\n")
	} else {
		info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.CashGamePayment, i18n.GetTranslator(paymentLexeme)(ctx), paymentType))
	}

	if name, isPresent := game.Name.Get(); isPresent {
		info.WriteString(fmt.Sprintf("%s %s: %s %s\n", icons.Sharp, i18n.GetTranslator(titleLexeme)(ctx), name, game.Number))
	} else {
		info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Sharp, i18n.GetTranslator(numberLexeme)(ctx), game.Number))
	}

	if game.Price > 0 {
		info.WriteString(fmt.Sprintf("%s %s: %d₽\n", icons.USD, i18n.GetTranslator(gameCostLexeme)(ctx), game.Price))
	}

	info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Place, i18n.GetTranslator(addressLexeme)(ctx), place.Address))
	info.WriteString(fmt.Sprintf("%s %s: %s\n", icons.Calendar, i18n.GetTranslator(dateTimeLexeme)(ctx), game.DateTime))
	info.WriteString(fmt.Sprintf("%s %s: %d/%d/%d", icons.NumberOfPlayers, i18n.GetTranslator(numberOfPlayersLexeme)(ctx), numberOfPlayers, numberOfLegioners, game.MaxPlayers))

	return info.String()
}
