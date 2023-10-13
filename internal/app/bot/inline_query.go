package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

var (
	registrationForAGameLexeme = i18n.Lexeme{
		Key:      "registration_for_a_game",
		FallBack: "Registration for a game",
	}
	thereAreNoGamePlayersLexeme = i18n.Lexeme{
		Key:      "there_are_no_game_players_lexeme",
		FallBack: "There are no game players lexeme",
	}
	zoyaLexeme = i18n.Lexeme{
		Key:      "zoya",
		FallBack: "Zoya",
	}
	inviteLexeme = i18n.Lexeme{
		Key:      "i_invite_you_to_register_for_the_game",
		FallBack: "I invite you to register for the game",
	}
)

func (b *Bot) handleInlineQuery(ctx context.Context, update *tgbotapi.Update) error {
	fn := func(ctx context.Context, update *tgbotapi.Update) ([]tgbotapi.Chattable, error) {
		inlineQueryID := update.InlineQuery.ID

		var results []interface{}

		inline := tgbotapi.InlineConfig{
			InlineQueryID: inlineQueryID,
			IsPersonal:    false,
			CacheTime:     5,
			Results:       results,
		}

		if update.InlineQuery.Query == "tip" {
			article := tgbotapi.NewInlineQueryResultArticle(uuid.NewString(), "Send tip-message", i18n.GetTranslator(registrationForAGameLexeme)(ctx))
			btn := tgbotapi.NewInlineKeyboardButtonURL(i18n.GetTranslator(zoyaLexeme)(ctx), "https://t.me/quiz_regbot")
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
				btn,
			})
			article.ReplyMarkup = &replyMarkup
			inline.Results = append(inline.Results, article)

			return []tgbotapi.Chattable{inline}, nil
		}

		games, err := b.gamesFacade.GetGames(ctx, true, true, false)
		if err != nil {
			return nil, fmt.Errorf("getting game error: %w", err)
		}

		if len(games) == 0 {
			return []tgbotapi.Chattable{inline}, nil
		}

		for _, game := range games {
			league, err := b.leaguesFacade.GetLeague(ctx, game.LeagueID)
			if err != nil {
				return nil, fmt.Errorf("getting league error: %w", err)
			}

			place, err := b.placesFacade.GetPlace(ctx, game.PlaceID)
			if err != nil {
				return nil, fmt.Errorf("getting place error: %w", err)
			}

			title := fmt.Sprintf("%s %s", game.DateTime, place.Name)
			article := tgbotapi.NewInlineQueryResultArticle(uuid.NewString(), title, i18n.GetTranslator(inviteLexeme)(ctx))
			article.ThumbURL = league.LogoLink
			article.ThumbHeight = 100
			article.ThumbWidth = 100

			var gamePlayers []model.GamePlayer
			gamePlayers, err = b.gamePlayersFacade.GetGamePlayersByGameID(ctx, game.ID)
			if err != nil {
				return nil, fmt.Errorf("getting players by game ID error: %w", err)
			}

			textBuilder := strings.Builder{}

			if len(gamePlayers) == 0 {
				textBuilder.WriteString(i18n.GetTranslator(thereAreNoGamePlayersLexeme)(ctx))
			}

			for i, gamePlayer := range gamePlayers {
				playerName := ""
				if userID, ok := gamePlayer.UserID.Get(); ok {
					var user model.User
					if user, err = b.usersFacade.GetUser(ctx, userID); err != nil {
						return nil, fmt.Errorf("getting user error: %w", err)
					}
					playerName = user.Name
				} else {
					playerName = "Лег"
				}

				ss := strings.Split(playerName, " ")
				name := ss[0]
				if i < len(gamePlayers)-1 {
					textBuilder.WriteString(fmt.Sprintf("%s, ", name))
				} else {
					textBuilder.WriteString(name)
				}
			}

			article.Description = textBuilder.String()

			btn := tgbotapi.NewInlineKeyboardButtonData(title, strconv.Itoa(int(game.ID)))
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
				btn,
			})
			article.ReplyMarkup = &replyMarkup

			results = append(results, article)
		}

		inline.Results = results
		return []tgbotapi.Chattable{inline}, nil
	}

	inlines, err := fn(ctx, update)
	if err != nil {
		return fmt.Errorf("preparing inline error: %w", err)
	}

	for _, inline := range inlines {
		if _, err = b.bot.Request(inline); err != nil {
			return fmt.Errorf("sending inline error: %w", err)
		}
	}

	return nil
}
