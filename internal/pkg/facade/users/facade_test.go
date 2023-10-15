package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFacade(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got := NewFacade(Config{})
		assert.Equal(t, got, &Facade{})
	})
}
