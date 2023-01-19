package telegramapi

import (
	"context"
	"fmt"
	"net"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	telegrampb "github.com/nikita5637/quiz-telegram/pkg/pb/telegram"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TelegramBot ...
type TelegramBot interface { // nolint:revive
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
}

// Implementation ...
type Implementation struct {
	bindAddr   string
	bot        TelegramBot
	grpcServer *grpc.Server

	telegrampb.UnimplementedMessageSenderServiceServer
}

// Config ...
type Config struct {
	BindAddr string
	Bot      TelegramBot
}

// New ...
func New(cfg Config) (*Implementation, error) {
	implementation := &Implementation{
		bindAddr: cfg.BindAddr,
		bot:      cfg.Bot,
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.ChainUnaryInterceptor(
		grpc_recovery.UnaryServerInterceptor(),
		logInterceptor,
	))

	s := grpc.NewServer(opts...)
	telegrampb.RegisterMessageSenderServiceServer(s, implementation)
	reflection.Register(s)

	implementation.grpcServer = s

	return implementation, nil
}

// ListenAndServe ...
func (i *Implementation) ListenAndServe(ctx context.Context) error {
	var err error
	var lis net.Listener

	lis, err = net.Listen("tcp", i.bindAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		err = i.grpcServer.Serve(lis)
		return
	}()
	if err != nil {
		return err
	}

	<-ctx.Done()

	i.grpcServer.GracefulStop()

	logger.Info(ctx, "telegram API gracefully stopped")
	return nil
}
