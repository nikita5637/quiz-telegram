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

	resp, err := b.registratorServiceClient.GetRegisteredGames(ctx, &registrator.GetRegisteredGamesRequest{
		Active: true,
	})
	if err != nil {
		return err
	}

	if len(resp.GetGames()) == 0 {
		_, err = b.bot.Request(inline)
		return err
	}

	for _, game := range resp.GetGames() {
		var gameResp *registrator.GetPlaceByIDResponse
		gameResp, err = b.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
			Id: game.GetPlaceId(),
		})
		if err != nil {
			return err
		}

		title := fmt.Sprintf("%s %s", model.DateTime(game.GetDate().AsTime()), gameResp.GetPlace().GetName())
		article := tgbotapi.NewInlineQueryResultArticle(uuid.NewString(), title, "Приглашаю зарегистрироваться на игру")
		if game.LeagueId == 1 {
			article.ThumbURL = "https://quizplease.ru//img/header-logo-white-2.png"
			article.ThumbHeight = 100
			article.ThumbWidth = 100
		} else if game.LeagueId == 2 {
			article.ThumbURL = "https://static.tildacdn.com/tild3931-3231-4334-b961-653062386133/_2.png"
			article.ThumbHeight = 100
			article.ThumbWidth = 100
		} else if game.LeagueId == 4 {
			article.ThumbURL = "https://sun9-14.userapi.com/impg/rokB7N7skTBT_4LELlQx6equHNPwMWPeGSQ5bQ/RykvYSWo5CA.jpg?size=500x500&quality=95&sign=40621a7260bd5eb2c20ae16111b537cc&type=album"
			article.ThumbHeight = 100
			article.ThumbWidth = 100
		}

		var playersResp *registrator.GetPlayersByGameIDResponse
		playersResp, err = b.registratorServiceClient.GetPlayersByGameID(ctx, &registrator.GetPlayersByGameIDRequest{
			GameId: game.GetId(),
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
				var playerResp *registrator.GetUserByIDResponse
				if playerResp, err = b.registratorServiceClient.GetUserByID(ctx, &registrator.GetUserByIDRequest{
					Id: player.GetUserId(),
				}); err != nil {
					return err
				}
				playerName = playerResp.GetUser().GetName()
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

		btn := tgbotapi.NewInlineKeyboardButtonData(title, strconv.Itoa(int(game.GetId())))
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
