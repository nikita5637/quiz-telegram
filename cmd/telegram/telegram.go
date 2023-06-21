package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"

	icsfilemanagerpb "github.com/nikita5637/quiz-ics-manager-api/pkg/pb/ics_file_manager"
	croupierpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/croupier"
	photomanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/photo_manager"
	"github.com/nikita5637/quiz-registrator-api/pkg/pb/registrator"
	telegram "github.com/nikita5637/quiz-telegram/internal/app/bot"
	"github.com/nikita5637/quiz-telegram/internal/app/reminder"
	"github.com/nikita5637/quiz-telegram/internal/app/telegramapi"
	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/elasticsearch"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/gamephotos"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/games"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/icsfiles"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/leagues"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/places"
	"github.com/nikita5637/quiz-telegram/internal/pkg/facade/users"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/middleware"
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
	flag.Parse()

	ctx := context.Background()

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		oscall := <-c
		logger.Infof(ctx, "system call recieved: %+v", oscall)
		cancel()
	}()

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
	registratorAPIClientConn, err := grpc.Dial(target, opts, grpc.WithChainUnaryInterceptor(
		middleware.LogInterceptor,
		middleware.TelegramClientIDInterceptor,
	))
	if err != nil {
		panic(err)
	}
	defer registratorAPIClientConn.Close()

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

	croupierServiceClient := croupierpb.NewServiceClient(registratorAPIClientConn)
	photographerServiceClient := photomanagerpb.NewServiceClient(registratorAPIClientConn)
	registratorServiceClient := registrator.NewRegistratorServiceClient(registratorAPIClientConn)
	icsFileManagerAPIServiceClient := icsfilemanagerpb.NewServiceClient(icsManagerAPIClientConn)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		leaguesFacadeConfig := leagues.Config{
			RegistratorServiceClient: registratorServiceClient,
		}
		leaguesFacade := leagues.NewFacade(leaguesFacadeConfig)

		placesFacadeConfig := places.Config{
			RegistratorServiceClient: registratorServiceClient,
		}
		placesFacade := places.NewFacade(placesFacadeConfig)

		gamePhotosFacadeConfig := gamephotos.Config{
			LeaguesFacade: leaguesFacade,
			PlacesFacade:  placesFacade,

			PhotographerServiceClient: photographerServiceClient,
			RegistratorServiceClient:  registratorServiceClient,
		}
		gamePhotosFacade := gamephotos.NewFacade(gamePhotosFacadeConfig)

		gamesFacadeConfig := games.Config{
			LeaguesFacade: leaguesFacade,
			PlacesFacade:  placesFacade,

			RegistratorServiceClient: registratorServiceClient,
		}
		gamesFacade := games.NewFacade(gamesFacadeConfig)

		icsFilesFacadeConfig := icsfiles.Config{
			ICSFileManagerAPIServiceClient: icsFileManagerAPIServiceClient,
		}
		icsFilesFacade := icsfiles.NewFacade(icsFilesFacadeConfig)

		usersFacadeConfig := users.Config{
			RegistratorServiceClient: registratorServiceClient,
		}
		usersFacade := users.NewFacade(usersFacadeConfig)

		telegramBotConfig := telegram.Config{
			Bot:              bot,
			GamePhotosFacade: gamePhotosFacade,
			GamesFacade:      gamesFacade,
			ICSFilesFacade:   icsFilesFacade,
			PlacesFacade:     placesFacade,
			UsersFacade:      usersFacade,

			CroupierServiceClient: croupierServiceClient,
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
