package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorCodes(t *testing.T) {
	assert.Equal(t, "\x1b[32mhello\x1b[0m", colorInfo("hello"), "they should be equal")
	assert.Equal(t, "\x1b[33mhello\x1b[0m", colorWarn("hello"), "they should be equal")
	assert.Equal(t, "\x1b[31mhello\x1b[0m", colorErr("hello"), "they should be equal")
}
