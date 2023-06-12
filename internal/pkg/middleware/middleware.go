package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	telegram_utils "github.com/nikita5637/quiz-telegram/utils/telegram"
	time_utils "github.com/nikita5637/quiz-telegram/utils/time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LogInterceptor ...
func LogInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time_utils.TimeNow()
	err := invoker(ctx, method, req, reply, cc, opts...)
	logger.Debugf(ctx, "Invoked RPC method=%s; Duration=%s; Error=%v", method, time.Since(start), err)
	return err
}

// TelegramClientIDInterceptor ...
func TelegramClientIDInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	telegramClientID := telegram_utils.ClientIDFromContext(ctx)
	if telegramClientID != 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, telegram_utils.TelegramClientID, strconv.FormatInt(telegramClientID, 10))
	}

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
