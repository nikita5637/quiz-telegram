package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	icsfilemanagerpb "github.com/nikita5637/quiz-ics-manager-api/pkg/pb/ics_file_manager"
	certificatemanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/certificate_manager"
	croupierpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/croupier"
	gamepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game"
	gameplayerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_player"
	gameresultmanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/game_result_manager"
	leaguepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/league"
	photomanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/photo_manager"
	placepb "github.com/nikita5637/quiz-registrator-api/pkg/pb/place"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	telegram "github.com/nikita5637/quiz-telegram/internal/app/bot"
	"github.com/nikita5637/quiz-telegram/internal/app/reminder"
	"github.com/nikita5637/quiz-telegram/internal/app/telegramapi"
	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/elasticsearch"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/certificates"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gamephotos"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gameplayers"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gameresults"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/icsfiles"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/leagues"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/places"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/users"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/middleware"
	"github.com/posener/ctxutil"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./config.toml", "path to config file")
}

func main() {
	ctx := ctxutil.Interrupt()
	flag.Parse()

	var err error
	err = config.ParseConfigFile(configPath)
	if err != nil {
		panic(err)
	}

	logsCombiner := &logger.Combiner{}
	logsCombiner = logsCombiner.WithWriter(os.Stdout)

	elasticLogsEnabled := config.GetValue("ElasticLogsEnabled").Bool()
	if elasticLogsEnabled {
		var elasticClient *elasticsearch.Client
		elasticClient, err = elasticsearch.New(elasticsearch.Config{
			ElasticAddress: config.GetElasticAddress(),
			ElasticIndex:   config.GetValue("ElasticIndex").String(),
		})
		if err != nil {
			panic(err)
		}

		logger.Info(ctx, "initialized elasticsearch client")
		logsCombiner = logsCombiner.WithWriter(elasticClient)
	}

	logLevel := config.GetLogLevel()
	logger.SetGlobalLogger(logger.NewLogger(logLevel, logsCombiner, zap.Fields(
		zap.String("module", "telegram"),
	)))
	logger.InfoKV(ctx, "initialized logger", "log level", logLevel)

	var bot *tgbotapi.BotAPI
	bot, err = tgbotapi.NewBotAPI(config.GetSecretValue(config.TelegramToken))
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	logger.Infof(ctx, "authorized on account '%s'", bot.Self.UserName)

	registratorAPIAddress := config.GetValue("RegistratorAPIAddress").String()
	registratorAPIPort := config.GetValue("RegistratorAPIPort").Uint16()

	opts := grpc.WithInsecure()
	target := fmt.Sprintf("%s:%d", registratorAPIAddress, registratorAPIPort)
	registratorAPIClientServiceConn, err := grpc.Dial(target, opts, grpc.WithChainUnaryInterceptor(
		middleware.LogInterceptor,
		middleware.ServiceNameInterceptor,
	))
	if err != nil {
		logger.Fatalf(ctx, "registratorAPIClient service conn dial error: %w", err.Error())
	}
	defer registratorAPIClientServiceConn.Close()

	registratorAPIClientUserConn, err := grpc.Dial(target, opts, grpc.WithChainUnaryInterceptor(
		middleware.LogInterceptor,
		middleware.TelegramClientIDInterceptor,
	))
	if err != nil {
		logger.Fatalf(ctx, "registratorAPIClient user conn dial error: %w", err.Error())
	}
	defer registratorAPIClientUserConn.Close()

	icsManagerAPIAddress := config.GetValue("ICSManagerAPIAddress").String()
	icsManagerAPIPort := config.GetValue("ICSManagerAPIPort").Uint16()

	target = fmt.Sprintf("%s:%d", icsManagerAPIAddress, icsManagerAPIPort)
	icsManagerAPIClientConn, err := grpc.Dial(target, opts, grpc.WithChainUnaryInterceptor(
		middleware.LogInterceptor,
		middleware.TelegramClientIDInterceptor,
	))
	if err != nil {
		panic(err)
	}
	defer icsManagerAPIClientConn.Close()

	certificateManagerServiceClient := certificatemanagerpb.NewServiceClient(registratorAPIClientUserConn)
	croupierServiceClient := croupierpb.NewServiceClient(registratorAPIClientUserConn)
	gamePlayerServiceClient := gameplayerpb.NewServiceClient(registratorAPIClientUserConn)
	gamePlayerRegistratorServiceClient := gameplayerpb.NewRegistratorServiceClient(registratorAPIClientUserConn)
	leagueServiceClient := leaguepb.NewServiceClient(registratorAPIClientUserConn)
	photographerServiceClient := photomanagerpb.NewServiceClient(registratorAPIClientUserConn)
	placeServiceClient := placepb.NewServiceClient(registratorAPIClientUserConn)
	gameServiceClient := gamepb.NewServiceClient(registratorAPIClientUserConn)
	gameRegistratorServiceClient := gamepb.NewRegistratorServiceClient(registratorAPIClientUserConn)
	gameResultManagerClient := gameresultmanagerpb.NewServiceClient(registratorAPIClientUserConn)
	userManagerServiceUserClient := usermanagerpb.NewServiceClient(registratorAPIClientUserConn)
	userManagerServiceServiceClient := usermanagerpb.NewServiceClient(registratorAPIClientServiceConn)
	icsFileManagerAPIServiceClient := icsfilemanagerpb.NewServiceClient(icsManagerAPIClientConn)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		leaguesFacadeConfig := leagues.Config{
			LeagueServiceClient: leagueServiceClient,
		}
		leaguesFacade := leagues.NewFacade(leaguesFacadeConfig)

		placesFacadeConfig := places.Config{
			PlaceServiceClient: placeServiceClient,
		}
		placesFacade := places.NewFacade(placesFacadeConfig)

		certificatesFacadeConfig := certificates.Config{
			CertificateManagerServiceClient: certificateManagerServiceClient,
		}
		certificatesFacade := certificates.New(certificatesFacadeConfig)

		gamePhotosFacadeConfig := gamephotos.Config{
			PhotographerServiceClient: photographerServiceClient,
		}
		gamePhotosFacade := gamephotos.New(gamePhotosFacadeConfig)

		gamePlayersFacadeConfig := gameplayers.Config{
			GamePlayerServiceClient:            gamePlayerServiceClient,
			GamePlayerRegistratorServiceClient: gamePlayerRegistratorServiceClient,
		}
		gamePlayersFacade := gameplayers.New(gamePlayersFacadeConfig)

		gamesFacadeConfig := games.Config{
			GamePlayersFacade: gamePlayersFacade,

			GameServiceClient:            gameServiceClient,
			GameRegistratorServiceClient: gameRegistratorServiceClient,
		}
		gamesFacade := games.New(gamesFacadeConfig)

		gameResultsFacadeConfig := gameresults.Config{
			GameResultManagerClient: gameResultManagerClient,
		}
		gameResultsFacade := gameresults.New(gameResultsFacadeConfig)

		icsFilesFacadeConfig := icsfiles.Config{
			ICSFileManagerAPIServiceClient: icsFileManagerAPIServiceClient,
		}
		icsFilesFacade := icsfiles.NewFacade(icsFilesFacadeConfig)

		usersFacadeConfig := users.Config{
			UserManagerServiceClient: userManagerServiceUserClient,
		}
		usersFacade := users.NewFacade(usersFacadeConfig)

		telegramBotConfig := telegram.Config{
			Bot:                bot,
			CertificatesFacade: certificatesFacade,
			GamesFacade:        gamesFacade,
			GameResultsFacade:  gameResultsFacade,
			GamePhotosFacade:   gamePhotosFacade,
			GamePlayersFacade:  gamePlayersFacade,
			ICSFilesFacade:     icsFilesFacade,
			LeaguesFacade:      leaguesFacade,
			PlacesFacade:       placesFacade,
			UsersFacade:        usersFacade,

			CroupierServiceClient:    croupierServiceClient,
			UserManagerServiceClient: userManagerServiceServiceClient,
		}

		tgBot, err2 := telegram.New(telegramBotConfig)
		if err2 != nil {
			return err2
		}

		logger.InfoKV(ctx, "initialized telegram bot", "registrator api address", registratorAPIAddress, "registrator api port", registratorAPIPort)
		return tgBot.Start(ctx)
	})

	g.Go(func() error {
		bindAddr := config.GetTelegramAPIBindAddress()
		telegramAPIConfig := telegramapi.Config{
			BindAddr: bindAddr,
			Bot:      bot,
		}

		telegramAPI, err2 := telegramapi.New(telegramAPIConfig)
		if err2 != nil {
			return err2
		}

		logger.InfoKV(ctx, "initialized telegram API", "bind address", bindAddr)
		return telegramAPI.ListenAndServe(ctx)
	})

	g.Go(func() error {
		rabbitMQConn, err2 := amqp.Dial(config.GetRabbitMQURL())
		if err2 != nil {
			return err2
		}
		defer rabbitMQConn.Close()

		rabbitMQChannel, err2 := rabbitMQConn.Channel()
		if err2 != nil {
			return err2
		}
		defer rabbitMQChannel.Close()

		gameReminderQueueName := config.GetValue("RabbitMQGameReminderQueueName").String()
		if gameReminderQueueName == "" {
			return errors.New("empty rabbit MQ game reminder queue name")
		}

		lotteryReminderQueueName := config.GetValue("RabbitMQLotteryReminderQueueName").String()
		if lotteryReminderQueueName == "" {
			return errors.New("empty rabbit MQ lottery reminder queue name")
		}

		reminderConfig := reminder.Config{
			Bot:                      bot,
			GameReminderQueueName:    gameReminderQueueName,
			LotteryReminderQueueName: lotteryReminderQueueName,
			RabbitMQChannel:          rabbitMQChannel,
			RegistratorAPIAddress:    registratorAPIAddress,
			RegistratorAPIPort:       registratorAPIPort,
		}
		reminder := reminder.New(reminderConfig)

		logger.InfoKV(ctx, "initialized reminder")
		return reminder.Start(ctx)
	})

	err = g.Wait()
	if err != nil {
		logger.Panic(ctx, err)
	}
}
