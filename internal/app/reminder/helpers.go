package reminder

import (
	"context"

	"github.com/nikita5637/quiz-telegram/internal/pkg/i18n"
)

func getTranslator(lexeme i18n.Lexeme) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return i18n.Translate(ctx, lexeme.Key, lexeme.FallBack)
	}
}
