//go:generate mockery --case underscore --name CertificatesFacade --with-expecter
//go:generate mockery --case underscore --name GamesFacade --with-expecter
//go:generate mockery --case underscore --name GamePhotosFacade --with-expecter
//go:generate mockery --case underscore --name GamePlayersFacade --with-expecter
//go:generate mockery --case underscore --name GameResultsFacade --with-expecter
//go:generate mockery --case underscore --name ICSFilesFacade --with-expecter
//go:generate mockery --case underscore --name LeaguesFacade --with-expecter
//go:generate mockery --case underscore --name PlacesFacade --with-expecter
//go:generate mockery --case underscore --name UsersFacade --with-expecter
//go:generate mockery --case underscore --name CroupierServiceClient --with-expecter
//go:generate mockery --case underscore --name UserManagerServiceClient --with-expecter
//go:generate mockery --case underscore --name TelegramBot --with-expecter

package bot

import (
	"context"
	"runtime/debug"
	"sync"

	croupierpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/croupier"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	telegrampb "github.com/nikita5637/quiz-telegram/pkg/pb/telegram"
	"google.golang.org/grpc"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CertificatesFacade ...
type CertificatesFacade interface {
	GetActiveCertificates(ctx context.Context) ([]model.Certificate, error)
}

// GamesFacade ...
type GamesFacade interface {
	GetGame(ctx context.Context, id int32) (model.Game, error)
	GetGames(ctx context.Context, registered, isInMaster, hasPassed bool) ([]model.Game, error)
	GetGamesByUserID(ctx context.Context, userID int32) ([]model.Game, error)
	SearchPassedAndRegisteredGames(ctx context.Context, page, pageSize uint64) ([]model.Game, uint64, error)
	RegisterGame(ctx context.Context, gameID int32) error
	UnregisterGame(ctx context.Context, gameID int32) error
	UpdatePayment(ctx context.Context, gameID, payment int32) error
}

// GamePhotosFacade ...
type GamePhotosFacade interface {
	GetPhotosByGameID(ctx context.Context, gameID int32) ([]string, error)
}

// GamePlayersFacade ...
type GamePlayersFacade interface {
	GetGamePlayersByGameID(ctx context.Context, gameID int32) ([]model.GamePlayer, error)
	RegisterPlayer(ctx context.Context, gamePlayer model.GamePlayer) error
	UnregisterPlayer(ctx context.Context, gamePlayer model.GamePlayer) error
	UpdatePlayerRegistration(ctx context.Context, gamePlayer model.GamePlayer) error
}

// GameResultsFacade ...
type GameResultsFacade interface {
	GetGameResultByGameID(ctx context.Context, gameID int32) (model.GameResult, error)
}

// ICSFilesFacade ...
type ICSFilesFacade interface {
	GetICSFileByGameID(ctx context.Context, gameID int32) (model.ICSFile, error)
}

// LeaguesFacade ...
type LeaguesFacade interface {
	GetLeague(ctx context.Context, id int32) (model.League, error)
}

// PlacesFacade ...
type PlacesFacade interface {
	GetPlace(ctx context.Context, id int32) (model.Place, error)
}

// UsersFacade ...
type UsersFacade interface {
	GetUser(ctx context.Context, userID int32) (model.User, error)
	UpdateUserBirthdate(ctx context.Context, userID int32, birthdate string) error
	UpdateUserEmail(ctx context.Context, userID int32, email string) error
	UpdateUserName(ctx context.Context, userID int32, name string) error
	UpdateUserPhone(ctx context.Context, userID int32, phone string) error
	UpdateUserSex(ctx context.Context, userID int32, sex model.Sex) error
	UpdateUserState(ctx context.Context, userID, state int32) error
}

// CroupierServiceClient ...
type CroupierServiceClient interface {
	GetLotteryStatus(ctx context.Context, in *croupierpb.GetLotteryStatusRequest, opts ...grpc.CallOption) (*croupierpb.GetLotteryStatusResponse, error)
	RegisterForLottery(ctx context.Context, in *croupierpb.RegisterForLotteryRequest, opts ...grpc.CallOption) (*croupierpb.RegisterForLotteryResponse, error)
}

// UserManagerServiceClient ...
type UserManagerServiceClient interface {
	CreateUser(ctx context.Context, in *usermanagerpb.CreateUserRequest, opts ...grpc.CallOption) (*usermanagerpb.User, error)
	GetUserByTelegramID(ctx context.Context, in *usermanagerpb.GetUserByTelegramIDRequest, opts ...grpc.CallOption) (*usermanagerpb.User, error)
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
	bot                TelegramBot // *tgbotapi.BotAPI
	certificatesFacade CertificatesFacade
	gamesFacade        GamesFacade
	gameResultsFacade  GameResultsFacade
	gamePhotosFacade   GamePhotosFacade
	gamePlayersFacade  GamePlayersFacade
	icsFilesFacade     ICSFilesFacade
	leaguesFacade      LeaguesFacade
	placesFacade       PlacesFacade
	usersFacade        UsersFacade

	croupierServiceClient    CroupierServiceClient
	userManagerServiceClient UserManagerServiceClient

	telegrampb.UnimplementedMessageSenderServiceServer
}

// Config ...
type Config struct {
	Bot                TelegramBot // *tgbotapi.BotAPI
	CertificatesFacade CertificatesFacade
	GamesFacade        GamesFacade
	GameResultsFacade  GameResultsFacade
	GamePhotosFacade   GamePhotosFacade
	GamePlayersFacade  GamePlayersFacade
	ICSFilesFacade     ICSFilesFacade
	LeaguesFacade      LeaguesFacade
	PlacesFacade       PlacesFacade
	UsersFacade        UsersFacade

	CroupierServiceClient    croupierpb.ServiceClient
	UserManagerServiceClient usermanagerpb.ServiceClient
}

// New ...
func New(cfg Config) (*Bot, error) {
	bot := &Bot{
		bot:                cfg.Bot,
		certificatesFacade: cfg.CertificatesFacade,
		gamesFacade:        cfg.GamesFacade,
		gameResultsFacade:  cfg.GameResultsFacade,
		gamePhotosFacade:   cfg.GamePhotosFacade,
		gamePlayersFacade:  cfg.GamePlayersFacade,
		icsFilesFacade:     cfg.ICSFilesFacade,
		leaguesFacade:      cfg.LeaguesFacade,
		placesFacade:       cfg.PlacesFacade,
		usersFacade:        cfg.UsersFacade,

		croupierServiceClient:    cfg.CroupierServiceClient,
		userManagerServiceClient: cfg.UserManagerServiceClient,
	}
	return bot, nil
}

// Start ...
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	wg := sync.WaitGroup{}
	go func(ctx context.Context) {
		for update := range updates {
			go func(ctx context.Context, update tgbotapi.Update) {
				wg.Add(1)
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						logger.ErrorKV(ctx, "panic recovered", "r", r, "update", update, "stack", string(debug.Stack()))
					}
				}()

				if update.CallbackQuery == nil && update.Message == nil && update.InlineQuery == nil {
					return
				}

				if err := b.handleUpdate(ctx, &update); err != nil {
					logger.Errorf(ctx, "update handling error: %s", err.Error())
				}
			}(ctx, update)
		}
	}(ctx)

	<-ctx.Done()

	b.bot.StopReceivingUpdates()
	wg.Wait()

	logger.Info(ctx, "telegram bot gracefully stopped")
	return nil
}
