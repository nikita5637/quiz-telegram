package i18n

import "context"

// GetTranslator ...
func GetTranslator(lexeme Lexeme) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return translate(ctx, lexeme.Key, lexeme.FallBack)
	}
}
