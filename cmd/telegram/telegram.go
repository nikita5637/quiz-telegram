package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	telegram "github.com/nikita5637/quiz-telegram/internal/app/bot"
	"github.com/nikita5637/quiz-telegram/internal/app/telegramapi"
	"github.com/nikita5637/quiz-telegram/internal/config"
	"github.com/nikita5637/quiz-telegram/internal/pkg/elasticsearch"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	"github.com/nikita5637/quiz-telegram/internal/pkg/request"
	"github.com/nikita5637/quiz-telegram/internal/pkg/storage"
	"go.uber.org/zap"

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

	db, err := storage.NewDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

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

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		requestStorage := storage.NewRequestStorage(db)

		requestsFacadeConfig := request.Config{
			RequestStorage: requestStorage,
		}
		requestsFacade := request.NewFacade(requestsFacadeConfig)

		registratorAPIAddress := config.GetValue("RegistratorAPIAddress").String()
		registratorAPIPort := config.GetValue("RegistratorAPIPort").Uint16()

		telegramBotConfig := telegram.Config{
			Bot: bot,

			RequestsFacade: requestsFacade,

			RegistratorAPIAddress: registratorAPIAddress,
			RegistratorAPIPort:    registratorAPIPort,
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

	err = g.Wait()
	if err != nil {
		logger.Panic(ctx, err)
	}
}
