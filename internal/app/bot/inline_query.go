package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

var (
	registrationForAGameLexeme = i18n.Lexeme{
		Key:      "registration_for_a_game",
		FallBack: "Registration for a game",
	}
	zoyaLexeme = i18n.Lexeme{
		Key:      "zoya",
		FallBack: "Zoya",
	}
)

// HandleInlineQuery ...
func (b *Bot) HandleInlineQuery(ctx context.Context, update *tgbotapi.Update) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-inline-query", "true")

	inlineQueryID := update.InlineQuery.ID

	var results []interface{}

	inline := tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		IsPersonal:    false,
		CacheTime:     5,
		Results:       results,
	}

	if update.InlineQuery.Query == "tip" {
		article := tgbotapi.NewInlineQueryResultArticle(uuid.NewString(), "Send tip-message", getTranslator(registrationForAGameLexeme)(ctx))
		btn := tgbotapi.NewInlineKeyboardButtonURL(getTranslator(zoyaLexeme)(ctx), "https://t.me/quiz_regbot")
		replyMarkup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
			btn,
		})
		article.ReplyMarkup = &replyMarkup
		inline.Results = append(inline.Results, article)

		_, err := b.bot.Request(inline)
		return err
	}

	games, err := b.gamesFacade.GetRegisteredGames(ctx, true)
	if err != nil {
		return err
	}

	if len(games) == 0 {
		_, err = b.bot.Request(inline)
		return err
	}

	for _, game := range games {
		title := fmt.Sprintf("%s %s", game.DateTime(), game.Place.Name)
		article := tgbotapi.NewInlineQueryResultArticle(uuid.NewString(), title, "Приглашаю зарегистрироваться на игру")
		article.ThumbURL = game.League.LogoLink
		article.ThumbHeight = 100
		article.ThumbWidth = 100

		var playersResp *registrator.GetPlayersByGameIDResponse
		playersResp, err = b.registratorServiceClient.GetPlayersByGameID(ctx, &registrator.GetPlayersByGameIDRequest{
			GameId: game.ID,
		})
		if err != nil {
			continue
		}

		textBuilder := strings.Builder{}

		if len(playersResp.GetPlayers()) == 0 {
			textBuilder.WriteString("Нет игроков")
		}

		for i, player := range playersResp.GetPlayers() {
			playerName := ""
			if player.GetUserId() > 0 {
				var user model.User
				if user, err = b.usersFacade.GetUserByID(ctx, player.GetUserId()); err != nil {
					return err
				}
				playerName = user.Name
			} else {
				playerName = "Лег"
			}

			ss := strings.Split(playerName, " ")
			name := ss[0]
			if i < len(playersResp.GetPlayers())-1 {
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
	_, err = b.bot.Request(inline)
	return err
}
