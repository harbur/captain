package captain // import "github.com/harbur/captain"

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPre(t *testing.T) {
	app := App{
		Pre: []string{"echo test"},
	}
	res := Pre(app)

	assert.Nil(t, res, "No error returned")
}

func TestPreFail(t *testing.T) {
	app := App{
		Pre: []string{"nonexistingCommand"},
	}
	res := Pre(app)
	assert.NotNil(t, res, "Error returned")
}
