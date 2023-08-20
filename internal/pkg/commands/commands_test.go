package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommands_Number(t *testing.T) {
	assert.Equal(t, CommandsNumber, Command(20))
}
