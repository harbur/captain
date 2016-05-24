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

func TestPost(t *testing.T) {
	app := App{
		Post: []string{"echo test"},
	}
	res := Post(app)

	assert.Nil(t, res, "No error returned")
}

func TestPostFail(t *testing.T) {
	app := App{
		Post: []string{"nonexistingCommand"},
	}
	res := Post(app)
	assert.NotNil(t, res, "Error returned")
}

func TestDownloadFile(t *testing.T) {
	res := downloadFile("/tmp/captain.html", "https://github.com/harbur/captain")
	assert.Nil(t, res, "captain")
}

func TestFindLastVersion(t *testing.T) {
	res := findLastVersion()
	assert.NotNil(t, res, "Last version exists")
}
