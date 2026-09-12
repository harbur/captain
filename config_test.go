package captain // import "github.com/harbur/captain"

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var basedir, _ = os.Getwd()

func TestConfigFiles(t *testing.T) {
	c := configFile("captain.yml")
	sl := "captain.yml"
	assert.Equal(t, sl, c, "Should return possible config files")
}

func TestReadConfig(t *testing.T) {
	c := readConfig(configFile(basedir + "/test/Simple/captain.yml"))
	assert.NotNil(t, c, "Should return configuration")
}

func TestNewConfig(t *testing.T) {
	info("cwd %s", basedir)
	c := NewConfig("", basedir+"/test/Simple/captain.yml", false)
	assert.NotNil(t, c, "Should return captain.yml configuration")
}

func TestNewConfigInferringValues(t *testing.T) {
	c := NewConfig("", basedir+"/test/noCaptainYML/captain.yml", false)
	assert.NotNil(t, c, "Should return infered configuration")
}

func TestFilterConfigEmpty(t *testing.T) {
	c := NewConfig("", basedir+"/test/Simple/captain.yml", false)
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")

	res := c.FilterConfig("")
	assert.True(t, res, "Should return true")
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")
}

func TestFilterConfigNonExistent(t *testing.T) {
	c := NewConfig("", basedir+"/test/Simple/captain.yml", false)
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")

	res := c.FilterConfig("nonexistent")
	assert.False(t, res, "Should return false")
	assert.Equal(t, 0, len(c.GetApps()), "Should return 0 apps")
}

func TestFilterConfigWeb(t *testing.T) {
	c := NewConfig("", basedir+"/test/Simple/captain.yml", false)
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")

	c.FilterConfig("web")
	assert.Equal(t, 1, len(c.GetApps()), "Should return 1 app")
	assert.Equal(t, "Dockerfile", c.GetApp("web").Build, "Should return web Build field")
}

func TestGetApp(t *testing.T) {
	c := NewConfig("", basedir+"/test/Simple/captain.yml", false)
	app := c.GetApp("web")
	assert.Equal(t, "harbur/test_web", app.Image, "Should return web image")
}

func TestConfigFilesDefault(t *testing.T) {
	c := configFile("")
	assert.Equal(t, "captain.yml", c, "Should default to captain.yml when no path is given")
}

func TestGetAppNonExistent(t *testing.T) {
	c := NewConfig("", basedir+"/test/Simple/captain.yml", false)
	app := c.GetApp("nonexistent")
	assert.Equal(t, App{}, app, "Should return an empty App when the name is not found")
}

// displaySyntaxError is written to format a *json.SyntaxError, even though the
// current caller (unmarshal) only ever feeds it YAML errors - exercise it
// directly with the error type it actually knows how to render.
func TestDisplaySyntaxErrorWithJSONSyntaxError(t *testing.T) {
	data := []byte("line one\nline two\n{ invalid json")
	var v interface{}
	jsonErr := json.Unmarshal(data, &v)

	syntaxErr, ok := jsonErr.(*json.SyntaxError)
	if !ok {
		t.Fatalf("expected json.Unmarshal to return a *json.SyntaxError, got %T", jsonErr)
	}

	res := displaySyntaxError(data, syntaxErr)
	assert.Error(t, res, "Should format a human-readable error")
	assert.Contains(t, res.Error(), "Error in line", "Should mention the offending line")
}

func TestDisplaySyntaxErrorWithNonSyntaxError(t *testing.T) {
	original := errors.New("some other kind of failure")
	res := displaySyntaxError([]byte("irrelevant"), original)
	assert.Equal(t, original, res, "Should return the original error untouched when it isn't a *json.SyntaxError")
}
