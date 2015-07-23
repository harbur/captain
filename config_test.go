package captain // import "github.com/harbur/captain"

import (
	"testing"

	"github.com/harbur/captain/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestConfigFiles(t *testing.T) {
	options.config = "captain.yml"
	c := configFile(options)
	sl := "captain.yml"
	assert.Equal(t, sl, c, "Should return possible config files")
}

func TestReadConfig(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := readConfig(configFile(options))
	assert.NotNil(t, c, "Should return configuration")
}

func TestNewConfig(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options, false)
	assert.NotNil(t, c, "Should return captain.yml configuration")
}

func TestNewConfigInferringValues(t *testing.T) {
	options.config = "test/noCaptainYML/captain.yml"
	c := NewConfig(options, false)
	assert.NotNil(t, c, "Should return infered configuration")
}

func TestFilterConfigEmpty(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options, false)
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")

	res := c.FilterConfig("")
	assert.True(t, res, "Should return true")
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")
}

func TestFilterConfigNonExistent(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options, false)
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")

	res := c.FilterConfig("nonexistent")
	assert.False(t, res, "Should return false")
	assert.Equal(t, 0, len(c.GetApps()), "Should return 0 apps")
}

func TestFilterConfigWeb(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options, false)
	assert.Equal(t, 2, len(c.GetApps()), "Should return 2 apps")

	c.FilterConfig("web")
	assert.Equal(t, 1, len(c.GetApps()), "Should return 1 app")
	assert.Equal(t, "Dockerfile", c.GetApp("web").Build, "Should return web Build field")
}

func TestGetApp(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options, false)
	app := c.GetApp("web")
	assert.Equal(t, "harbur/test_web", app.Image, "Should return web image")
}
