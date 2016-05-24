package captain // import "github.com/harbur/captain"

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test Fixtures
var validApp = App{
	Image: "",
	Pre:   []string{"echo running pre"},
	Post:  []string{"echo running post"},
}

var invalidApp = App{
	Pre:  []string{"nonexistingPreCommand"},
	Post: []string{"nonexistingPostCommand"},
}

// Pre Command
func TestPre(t *testing.T) {
	res := Pre(validApp)
	assert.Nil(t, res, "No error returned")
}

func TestPreFail(t *testing.T) {
	res := Pre(invalidApp)
	assert.NotNil(t, res, "Error returned")
}

// Post Command
func TestPost(t *testing.T) {
	res := Post(validApp)
	assert.Nil(t, res, "No error returned")
}

func TestPostFail(t *testing.T) {
	res := Post(invalidApp)
	assert.NotNil(t, res, "Error returned")
}

// Build Command
func TestBuild(t *testing.T) {
	var testConfig = readConfig(configFile("test/Simple/captain.yml"))

	var buildOpts = BuildOptions{
		Config: testConfig,
	}

	Build(buildOpts)
}

// SelfUpdate Command
func TestSelfUpdate(t *testing.T) {
	SelfUpdate()
}

func TestDownloadFile(t *testing.T) {
	res := downloadFile("/tmp/captain.html", "https://github.com/harbur/captain")
	assert.Nil(t, res, "captain")
}

func TestFindLastVersion(t *testing.T) {
	res := findLastVersion()
	assert.NotNil(t, res, "Last version exists")
}
