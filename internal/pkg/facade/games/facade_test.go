//go:generate mockery --case underscore --name RegistratorServiceClient --with-expecter

package games

import (
	"testing"

	"github.com/nikita5637/quiz-telegram/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewFacade(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		f := NewFacade(Config{
			RegistratorServiceClient: nil,
		})

		assert.Equal(t, &Facade{
			leagueCache: make(map[int32]model.League),
			placeCache:  make(map[int32]model.Place),
		}, f)
	})
}
