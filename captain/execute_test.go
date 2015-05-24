package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	res := execute("echo", "testing")
	assert.Equal(t, nil, res, "it should execute without errors")
}

func TestOneliner(t *testing.T) {
	res := oneliner("echo", "testing")
	assert.Equal(t, "testing", res, "it should return the trimmed result")
}

func TestOnelinerTrimmed(t *testing.T) {
	res := oneliner("echo", "testing with spaces  ")
	assert.Equal(t, "testing with spaces", res, "it should return the trimmed result")
}
