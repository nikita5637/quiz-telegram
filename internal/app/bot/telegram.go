//go:generate mockery --case underscore --name GamesFacade --with-expecter
//go:generate mockery --case underscore --name GamePhotosFacade --with-expecter
//go:generate mockery --case underscore --name CroupierServiceClient --with-expecter
//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter
//go:generate mockery --case underscore --name TelegramBot --with-expecter

package bot

import (
	"context"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	telegrampb "github.com/nikita5637/quiz-telegram/pkg/pb/telegram"
	uuid_utils "github.com/nikita5637/quiz-telegram/utils/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
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
}

// GamePhotosFacade ...
type GamePhotosFacade interface {
	GetGamesWithPhotos(ctx context.Context, limit, offset uint32) ([]model.Game, uint32, error)
	GetPhotosByGameID(ctx context.Context, gameID int32) ([]string, error)
}

// CroupierServiceClient ...
type CroupierServiceClient interface {
	registrator.CroupierServiceClient
}

// RegistratorServiceClient ...
type RegistratorServiceClient interface {
	// CreateUser creates new user
	CreateUser(ctx context.Context, in *registrator.CreateUserRequest, opts ...grpc.CallOption) (*registrator.CreateUserResponse, error)
	// GetPlaceByID returns place by place ID
	GetPlaceByID(ctx context.Context, in *registrator.GetPlaceByIDRequest, opts ...grpc.CallOption) (*registrator.GetPlaceByIDResponse, error)
	// GetPlayersByGameID returns list of players by game ID
	GetPlayersByGameID(ctx context.Context, in *registrator.GetPlayersByGameIDRequest, opts ...grpc.CallOption) (*registrator.GetPlayersByGameIDResponse, error)
	// GetUserByID returns user by user ID
	GetUserByID(ctx context.Context, in *registrator.GetUserByIDRequest, opts ...grpc.CallOption) (*registrator.GetUserByIDResponse, error)
	// GetUserByTelegramID returns a user by telegram ID
	GetUserByTelegramID(ctx context.Context, in *registrator.GetUserByTelegramIDRequest, opts ...grpc.CallOption) (*registrator.GetUserByTelegramIDResponse, error)
	// RegisterGame registers game
	RegisterGame(ctx context.Context, in *registrator.RegisterGameRequest, opts ...grpc.CallOption) (*registrator.RegisterGameResponse, error)
	// RegisterPlayer registers player for a game
	RegisterPlayer(ctx context.Context, in *registrator.RegisterPlayerRequest, opts ...grpc.CallOption) (*registrator.RegisterPlayerResponse, error)
	// UnregisterGame unregisters game
	UnregisterGame(ctx context.Context, in *registrator.UnregisterGameRequest, opts ...grpc.CallOption) (*registrator.UnregisterGameResponse, error)
	// UnregisterPlayer unregisters player
	UnregisterPlayer(ctx context.Context, in *registrator.UnregisterPlayerRequest, opts ...grpc.CallOption) (*registrator.UnregisterPlayerResponse, error)
	// UpdateUserEmail updates a user email
	UpdateUserEmail(ctx context.Context, in *registrator.UpdateUserEmailRequest, opts ...grpc.CallOption) (*registrator.UpdateUserEmailResponse, error)
	// UpdateUserName updates a user's name
	UpdateUserName(ctx context.Context, in *registrator.UpdateUserNameRequest, opts ...grpc.CallOption) (*registrator.UpdateUserNameResponse, error)
	// UpdateUserPhone updates a user's phone
	UpdateUserPhone(ctx context.Context, in *registrator.UpdateUserPhoneRequest, opts ...grpc.CallOption) (*registrator.UpdateUserPhoneResponse, error)
	// UpdateUserState updates a user's state
	UpdateUserState(ctx context.Context, in *registrator.UpdateUserStateRequest, opts ...grpc.CallOption) (*registrator.UpdateUserStateResponse, error)
	// UpdatePayment updates payment
	UpdatePayment(ctx context.Context, in *registrator.UpdatePaymentRequest, opts ...grpc.CallOption) (*registrator.UpdatePaymentResponse, error)
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
	bot              TelegramBot // *tgbotapi.BotAPI
	gamesFacade      GamesFacade
	gamePhotosFacade GamePhotosFacade

	croupierServiceClient    CroupierServiceClient
	registratorServiceClient RegistratorServiceClient

	telegrampb.UnimplementedMessageSenderServiceServer
}

// Config ...
type Config struct {
	Bot              TelegramBot // *tgbotapi.BotAPI
	GamesFacade      GamesFacade
	GamePhotosFacade GamePhotosFacade

	CroupierServiceClient    registrator.CroupierServiceClient
	RegistratorServiceClient registrator.RegistratorServiceClient
}

// New ...
func New(cfg Config) (*Bot, error) {
	bot := &Bot{
		bot:              cfg.Bot,
		gamesFacade:      cfg.GamesFacade,
		gamePhotosFacade: cfg.GamePhotosFacade,

		croupierServiceClient:    cfg.CroupierServiceClient,
		registratorServiceClient: cfg.RegistratorServiceClient,
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

				groupUUID := uuid.New().String()
				ctx = uuid_utils.NewContextWithGroupUUID(ctx, groupUUID)

				if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
					if err := b.HandleCallbackQuery(ctx, &update); err != nil {
						logger.Errorf(ctx, "callback query handle error: %s", err)
						clientID := update.CallbackQuery.From.ID
						responseMessage := tgbotapi.NewMessage(clientID, getTranslator(somethingWentWrongLexeme)(ctx))
						if s, ok := status.FromError(err); ok {
							if s.Code() == codes.PermissionDenied {
								responseMessage = tgbotapi.NewMessage(clientID, getTranslator(permissionDeniedLexeme)(ctx))
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
