package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorCodes(t *testing.T) {
	assert.Equal(t, info("hello"), "\x1b[34mhello\x1b[0m", "they should be equal")
	assert.Equal(t, warn("hello"), "\x1b[33mhello\x1b[0m", "they should be equal")
	assert.Equal(t, err("hello"), "\x1b[31mhello\x1b[0m", "they should be equal")
}
