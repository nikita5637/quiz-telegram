package games

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFacade(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		f := NewFacade(Config{
			LeaguesFacade: nil,
			PlacesFacade:  nil,

			RegistratorServiceClient: nil,
		})

		assert.Equal(t, &Facade{}, f)
	})
}
