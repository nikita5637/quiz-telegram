package gamephotos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFacade(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		f := New(Config{
			PhotographerServiceClient: nil,
		})

		assert.Equal(t, &Facade{}, f)
	})
}
