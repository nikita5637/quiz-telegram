package reminder

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	reminder "github.com/nikita5637/quiz-registrator-api/pkg/reminder"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/icons"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

var (
	placeLexeme = i18n.Lexeme{
		Key:      "place",
		FallBack: "Place",
	}
	remindThatThereIsAGameTodayLexeme = i18n.Lexeme{
		Key:      "remind_that_there_is_a_game_today",
		FallBack: "Remind that there is a game today",
	}
	timeLexeme = i18n.Lexeme{
		Key:      "time",
		FallBack: "Time",
	}
)

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	registrator.RegistratorServiceClient
}

// TelegramBot ...
type TelegramBot interface { // nolint:revive
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
}

// Reminder ...
type Reminder struct {
	bot                      TelegramBot // *tgbotapi.BotAPI
	gameReminderQueueName    string
	rabbitMQChanel           *amqp.Channel
	registratorAPIAddress    string
	registratorAPIPort       uint16
	registratorServiceClient RegistratorServiceClient
}

// Config ...
type Config struct {
	Bot                   TelegramBot // *tgbotapi.BotAPI
	GameReminderQueueName string
	RabbitMQChannel       *amqp.Channel
	RegistratorAPIAddress string
	RegistratorAPIPort    uint16
}

// New ...
func New(cfg Config) *Reminder {
	return &Reminder{
		bot:                   cfg.Bot,
		gameReminderQueueName: cfg.GameReminderQueueName,
		rabbitMQChanel:        cfg.RabbitMQChannel,
		registratorAPIAddress: cfg.RegistratorAPIAddress,
		registratorAPIPort:    cfg.RegistratorAPIPort,
	}
}

// Start ...
func (r *Reminder) Start(ctx context.Context) error {
	opts := grpc.WithInsecure()
	target := fmt.Sprintf("%s:%d", r.registratorAPIAddress, r.registratorAPIPort)
	cc, err := grpc.Dial(target, opts, grpc.WithChainUnaryInterceptor(
		moduleNameInterceptor,
	))
	if err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}

	r.registratorServiceClient = registrator.NewRegistratorServiceClient(cc)

	gameReminderQueue, err := r.rabbitMQChanel.QueueDeclare(
		r.gameReminderQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	gameReminderMessages, err := r.rabbitMQChanel.Consume(
		gameReminderQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		go func(ctx context.Context) {
			for d := range gameReminderMessages {
				logger.Debugf(ctx, "accepted new game remind message: %s", d.Body)

				gameRemind := &reminder.Game{}
				err := json.Unmarshal(d.Body, gameRemind)
				if err != nil {
					logger.Errorf(ctx, "get game remind error: %s", err.Error())
					continue
				}

				gameResp, err := r.registratorServiceClient.GetGameByID(ctx, &registrator.GetGameByIDRequest{
					GameId: gameRemind.GameID,
				})
				if err != nil {
					logger.Errorf(ctx, "get game by ID error: %s", err.Error())
					continue
				}

				placeResp, err := r.registratorServiceClient.GetPlaceByID(ctx, &registrator.GetPlaceByIDRequest{
					Id: gameResp.GetGame().GetPlaceId(),
				})
				if err != nil {
					logger.Errorf(ctx, "get place by ID error: %s", err.Error())
					continue
				}

				text := fmt.Sprintf("%s %s\n", icons.Note, getTranslator(remindThatThereIsAGameTodayLexeme)(ctx))
				text += fmt.Sprintf("%s %s: %s\n", icons.Time, getTranslator(timeLexeme)(ctx), model.DateTime(gameResp.GetGame().GetDate().AsTime()).Time())
				text += fmt.Sprintf("%s %s: %s\n", icons.Place, getTranslator(placeLexeme)(ctx), placeResp.GetPlace().GetAddress())

				for _, playerID := range gameRemind.PlayerIDs {
					resp, err := r.registratorServiceClient.GetUserByID(ctx, &registrator.GetUserByIDRequest{
						Id: playerID,
					})
					if err != nil {
						logger.Errorf(ctx, "get user by ID error: %s", err.Error())
						continue
					}

					textMessage := tgbotapi.NewMessage(resp.GetUser().GetTelegramId(), text)
					_, err = r.bot.Send(textMessage)
					if err != nil {
						logger.Errorf(ctx, "send game remind text message error: %s", err.Error())
						continue
					}

					venueMessage := tgbotapi.NewVenue(resp.GetUser().GetTelegramId(),
						placeResp.GetPlace().GetName(),
						placeResp.GetPlace().GetAddress(),
						float64(placeResp.GetPlace().GetLatitude()),
						float64(placeResp.GetPlace().GetLongitude()),
					)
					_, err = r.bot.Request(venueMessage)
					if err != nil {
						logger.Errorf(ctx, "send game remind venue message error: %s", err.Error())
						continue
					}
				}
			}
		}(ctx)
	}(ctx)

	<-ctx.Done()

	logger.Info(ctx, "reminder gracefully stopped")
	return nil
}
