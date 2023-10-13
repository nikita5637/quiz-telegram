package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	usermanagerpb "github.com/nikita5637/quiz-registrator-api/pkg/pb/user_manager"
	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
	"github.com/nikita5637/quiz-telegram/internal/pkg/logger"
	telegramutils "github.com/nikita5637/quiz-telegram/utils/telegram"
	userutils "github.com/nikita5637/quiz-telegram/utils/user"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	permissionDeniedLexeme = i18n.Lexeme{
		Key:      "permission_denied",
		FallBack: "Permission denied",
	}
	somethingWentWrongLexeme = i18n.Lexeme{
		Key:      "something_went_wrong",
		FallBack: "Something went wrong",
	}
	welcomeMessageLexeme = i18n.Lexeme{
		Key:      "welcome_message",
		FallBack: "Welcome message",
	}
)

// handleUpdate is a main handle function
func (b *Bot) handleUpdate(ctx context.Context, update *tgbotapi.Update) error {
	if update.Message != nil && !update.FromChat().IsPrivate() {
		logger.Debug(ctx, "not private message has been skipped")
		return nil
	}

	clientID := getClientIDFromUpdate(update)
	ctx = telegramutils.NewContextWithClientID(ctx, clientID)

	needToCreateUser := false
	pbUser, err := b.userManagerServiceClient.GetUserByTelegramID(ctx, &usermanagerpb.GetUserByTelegramIDRequest{
		TelegramId: clientID,
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() != codes.NotFound {
			return fmt.Errorf("getting user by telegram ID error: %w", err)
		}

		if st.Code() == codes.NotFound {
			if ei := getErrorInfoFromStatus(st); ei != nil {
				if ei.Reason != "USER_NOT_FOUND" {
					return fmt.Errorf("getting user by telegram ID error: %w", err)
				}
				needToCreateUser = true
			}
		}
	}

	name := getFirstNameFromUpdate(update)
	if name == "" {
		name = getUserNameFromUpdate(update)
	}

	if needToCreateUser {
		// TODO fix panix when creating new user when press button in chat
		if fromChat := update.FromChat(); fromChat != nil {
			if !fromChat.IsPrivate() {
				return nil
			}
		}

		if pbUser, err = b.userManagerServiceClient.CreateUser(ctx, &usermanagerpb.CreateUserRequest{
			User: &usermanagerpb.User{
				Name:       name,
				TelegramId: clientID,
				State:      usermanagerpb.UserState_USER_STATE_WELCOME,
			},
		}); err != nil {
			return fmt.Errorf("user creation error: %w", err)
		}
	}

	modelUser := convertProtoUserToModelUser(pbUser)
	if needToCreateUser {
		logger.InfoKV(ctx, "new user has been created", "user", modelUser)

		welcomeMessage := welcomeMessage(ctx, clientID, name)
		if _, err = b.bot.Send(welcomeMessage); err != nil {
			return err
		}

		return nil
	}

	logger.DebugKV(ctx, "user has been found by telegram ID", "user", modelUser)

	ctx = userutils.NewContextWithUser(ctx, &modelUser)

	var updateHandler func(ctx context.Context, update *tgbotapi.Update) error
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		if modelUser.State == int32(usermanagerpb.UserState_USER_STATE_CHANGING_BANNED) {
			permissionDeniedCallback := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(permissionDeniedLexeme)(ctx))
			if _, err := b.bot.Request(permissionDeniedCallback); err != nil {
				logger.Errorf(ctx, "sending permission denied callback error: %s", err.Error())
			}

			logger.InfoKV(ctx, "permission denied for user", "user", modelUser, "handler", "callback query")
			return nil
		}

		updateHandler = func(ctx context.Context, update *tgbotapi.Update) error {
			if err := b.handleCallbackQuery(ctx, update); err != nil {
				somethingWentWrongCallback := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(somethingWentWrongLexeme)(ctx))
				if _, err2 := b.bot.Request(somethingWentWrongCallback); err2 != nil {
					logger.Errorf(ctx, "sending something went wrong callback error: %s", err2.Error())
				}

				return fmt.Errorf("callback query handling error: %w", err)
			}

			return nil
		}
	} else if update.CallbackQuery != nil && update.CallbackQuery.InlineMessageID != "" {
		if modelUser.State == int32(usermanagerpb.UserState_USER_STATE_CHANGING_BANNED) {
			permissionDeniedCallback := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(permissionDeniedLexeme)(ctx))
			if _, err := b.bot.Request(permissionDeniedCallback); err != nil {
				logger.Errorf(ctx, "sending permission denied callback error: %s", err.Error())
			}

			logger.InfoKV(ctx, "permission denied for user", "user", modelUser, "handler", "inline message")
			return nil
		}

		updateHandler = func(ctx context.Context, update *tgbotapi.Update) error {
			if err := b.handleInlineMessage(ctx, update); err != nil {
				somethingWentWrongCallback := tgbotapi.NewCallback(update.CallbackQuery.ID, i18n.GetTranslator(somethingWentWrongLexeme)(ctx))
				if _, err2 := b.bot.Request(somethingWentWrongCallback); err2 != nil {
					logger.Errorf(ctx, "sending something went wrong callback error: %s", err2.Error())
				}

				return fmt.Errorf("inline message handling error: %w", err)
			}

			return nil
		}
	} else if update.InlineQuery != nil {
		if modelUser.State == int32(usermanagerpb.UserState_USER_STATE_CHANGING_BANNED) {
			logger.InfoKV(ctx, "permission denied for user", "user", modelUser, "handler", "inline query")
			return nil
		}

		updateHandler = func(ctx context.Context, update *tgbotapi.Update) error {
			if err := b.handleInlineQuery(ctx, update); err != nil {
				return fmt.Errorf("inline query handling error: %w", err)
			}

			return nil
		}
	} else if update.Message != nil {
		if modelUser.State == int32(usermanagerpb.UserState_USER_STATE_CHANGING_BANNED) {
			permissionDeniedMessage := tgbotapi.NewMessage(modelUser.TelegramID, i18n.GetTranslator(permissionDeniedLexeme)(ctx))
			if _, err := b.bot.Send(permissionDeniedMessage); err != nil {
				logger.Errorf(ctx, "sending permission denied message error: %s", err.Error())
			}

			logger.InfoKV(ctx, "permission denied for user", "user", modelUser)
			return nil
		}

		updateHandler = func(ctx context.Context, update *tgbotapi.Update) error {
			if err := b.handleMessage(ctx, update); err != nil {
				somethingWentWrongMessage := tgbotapi.NewMessage(modelUser.TelegramID, i18n.GetTranslator(somethingWentWrongLexeme)(ctx))
				if _, err2 := b.bot.Send(somethingWentWrongMessage); err2 != nil {
					logger.Errorf(ctx, "sending something went wrong message error: %s", err2.Error())
				}

				return fmt.Errorf("message handling error: %w", err)
			}

			return nil
		}
	}
	if updateHandler != nil {
		return updateHandler(ctx, update)
	}

	return nil
}

func getClientIDFromUpdate(update *tgbotapi.Update) int64 {
	if sentFrom := update.SentFrom(); sentFrom != nil {
		return sentFrom.ID
	}

	return 0
}

func getErrorInfoFromStatus(st *status.Status) *errdetails.ErrorInfo {
	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			return d
		}
	}

	return nil
}

func getFirstNameFromUpdate(update *tgbotapi.Update) string {
	if sentFrom := update.SentFrom(); sentFrom != nil {
		return sentFrom.FirstName
	}

	return ""
}

func getUserNameFromUpdate(update *tgbotapi.Update) string {
	if sentFrom := update.SentFrom(); sentFrom != nil {
		return sentFrom.UserName
	}

	return ""
}

func welcomeMessage(ctx context.Context, clientID int64, name string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(clientID, fmt.Sprintf(i18n.GetTranslator(welcomeMessageLexeme)(ctx), name))
}
