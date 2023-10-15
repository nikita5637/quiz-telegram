package leagues

import (
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewFacade(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got := NewFacade(Config{})
		assert.Equal(t, got, &Facade{
			leaguesCache: make(map[int32]model.League, 0),
		})
	})
}
