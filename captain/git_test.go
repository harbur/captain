package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitIsDirty(t *testing.T) {
	assert.Equal(t, false, isDirty(), "they should be equal")
}
