//go:generate mockery --case underscore --name GamesFacade --with-expecter
//go:generate mockery --case underscore --name GamePhotosFacade --with-expecter
//go:generate mockery --case underscore --name GamePlayersFacade --with-expecter
//go:generate mockery --case underscore --name ICSFilesFacade --with-expecter
//go:generate mockery --case underscore --name PlacesFacade --with-expecter
//go:generate mockery --case underscore --name UsersFacade --with-expecter
//go:generate mockery --case underscore --name CroupierServiceClient --with-expecter
//go:generate mockery --case underscore --name TelegramBot --with-expecter

package bot

import (
	"context"
	"runtime/debug"

	croupierpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/croupier"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	telegrampb "github.com/nikita5637/quiz-telegram/pkg/pb/telegram"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GamesFacade ...
type GamesFacade interface {
	GetGameByID(ctx context.Context, id int32) (model.Game, error)
	GetGames(ctx context.Context, active bool) ([]model.Game, error)
	GetRegisteredGames(ctx context.Context, active bool) ([]model.Game, error)
	GetUserGames(ctx context.Context, active bool, userID int32) ([]model.Game, error)
	RegisterGame(ctx context.Context, gameID int32) (int32, error)
	UnregisterGame(ctx context.Context, gameID int32) (int32, error)
	UpdatePayment(ctx context.Context, gameID, payment int32) error
}

// GamePhotosFacade ...
type GamePhotosFacade interface {
	GetGamesWithPhotos(ctx context.Context, limit, offset uint32) ([]model.Game, uint32, error)
	GetPhotosByGameID(ctx context.Context, gameID int32) ([]string, error)
}

// GamePlayersFacade ...
type GamePlayersFacade interface {
	GetGamePlayersByGameID(ctx context.Context, gameID int32) ([]model.GamePlayer, error)
	RegisterPlayer(ctx context.Context, gamePlayer model.GamePlayer) error
	UnregisterPlayer(ctx context.Context, gamePlayer model.GamePlayer) error
	UpdatePlayerRegistration(ctx context.Context, gamePlayer model.GamePlayer) error
}

// ICSFilesFacade ...
type ICSFilesFacade interface {
	GetICSFileByGameID(ctx context.Context, gameID int32) (model.ICSFile, error)
}

// PlacesFacade ...
type PlacesFacade interface {
	GetPlaceByID(ctx context.Context, id int32) (model.Place, error)
}

// UsersFacade ...
type UsersFacade interface {
	CreateUser(ctx context.Context, name string, telegramID int64, state int32) (int32, error)
	GetUserByID(ctx context.Context, userID int32) (model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (model.User, error)
	UpdateUserBirthdate(ctx context.Context, userID int32, birthdate string) error
	UpdateUserEmail(ctx context.Context, userID int32, email string) error
	UpdateUserName(ctx context.Context, userID int32, name string) error
	UpdateUserPhone(ctx context.Context, userID int32, phone string) error
	UpdateUserState(ctx context.Context, userID, state int32) error
	UpdateUserSex(ctx context.Context, userID int32, sex model.Sex) error
}

// CroupierServiceClient ...
type CroupierServiceClient interface {
	croupierpb.ServiceClient
}

// TelegramBot ...
type TelegramBot interface { // nolint:revive
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
}

// Bot ...
type Bot struct {
	bot               TelegramBot // *tgbotapi.BotAPI
	gamesFacade       GamesFacade
	gamePhotosFacade  GamePhotosFacade
	gamePlayersFacade GamePlayersFacade
	icsFilesFacade    ICSFilesFacade
	placesFacade      PlacesFacade
	usersFacade       UsersFacade

	croupierServiceClient CroupierServiceClient

	telegrampb.UnimplementedMessageSenderServiceServer
}

// Config ...
type Config struct {
	Bot               TelegramBot // *tgbotapi.BotAPI
	GamesFacade       GamesFacade
	GamePhotosFacade  GamePhotosFacade
	GamePlayersFacade GamePlayersFacade
	ICSFilesFacade    ICSFilesFacade
	PlacesFacade      PlacesFacade
	UsersFacade       UsersFacade

	CroupierServiceClient croupierpb.ServiceClient
}

// New ...
func New(cfg Config) (*Bot, error) {
	bot := &Bot{
		bot:               cfg.Bot,
		gamesFacade:       cfg.GamesFacade,
		gamePhotosFacade:  cfg.GamePhotosFacade,
		gamePlayersFacade: cfg.GamePlayersFacade,
		icsFilesFacade:    cfg.ICSFilesFacade,
		placesFacade:      cfg.PlacesFacade,
		usersFacade:       cfg.UsersFacade,

		croupierServiceClient: cfg.CroupierServiceClient,
	}
	return bot, nil
}

// Start ...
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	go func(ctx context.Context) {
		for update := range updates {
			go func(ctx context.Context, update tgbotapi.Update) {
				defer func() {
					if r := recover(); r != nil {
						logger.ErrorKV(ctx, "panic recovered", "r", r, "update", update, "stack", string(debug.Stack()))
					}
				}()

				if update.CallbackQuery == nil && update.Message == nil && update.InlineQuery == nil {
					return
				}

				if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
					if err := b.HandleCallbackQuery(ctx, &update); err != nil {
						logger.Errorf(ctx, "callback query handle error: %s", err)
						clientID := update.CallbackQuery.From.ID
						responseMessage := tgbotapi.NewMessage(clientID, i18n.GetTranslator(somethingWentWrongLexeme)(ctx))
						if s, ok := status.FromError(err); ok {
							if s.Code() == codes.PermissionDenied {
								responseMessage = tgbotapi.NewMessage(clientID, i18n.GetTranslator(permissionDeniedLexeme)(ctx))
							}
							if s.Code() == codes.NotFound {
								for _, detail := range s.Details() {
									switch t := detail.(type) {
									case *errdetails.LocalizedMessage:
										responseMessage = tgbotapi.NewMessage(clientID, t.GetMessage())
									}
								}
							}
						}

						if _, err := b.bot.Send(responseMessage); err != nil {
							logger.Errorf(ctx, "error while send message: %s", err)
						}
					}
				} else if update.CallbackQuery != nil && update.CallbackQuery.InlineMessageID != "" {
					var err2 error
					err2 = b.HandleInlineMessage(ctx, &update)
					if err2 != nil {
						logger.Errorf(ctx, "inline message handle error: %s", err2)
					}
				} else if update.InlineQuery != nil {
					var err2 error
					err2 = b.HandleInlineQuery(ctx, &update)
					if err2 != nil {
						logger.Errorf(ctx, "inline query handle error: %s", err2)
					}
				} else if update.Message != nil {
					if err := b.HandleMessage(ctx, &update); err != nil {
						logger.Errorf(ctx, "handle message error: %s", err)
					}
				}
			}(ctx, update)
		}
	}(ctx)

	<-ctx.Done()

	b.bot.StopReceivingUpdates()

	logger.Info(ctx, "telegram bot gracefully stopped")
	return nil
}
