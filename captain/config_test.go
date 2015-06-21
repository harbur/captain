package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	c := NewConfig(options,false)
	assert.NotNil(t, c, "Should return captain.yml configuration")
}

func TestNewConfigInferringValues(t *testing.T) {
	options.config = "test/noCaptainYML/captain.yml"
	c := NewConfig(options,false)
	assert.NotNil(t, c, "Should return infered configuration")
}

func TestGetImageNames(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options,false)
	expected := map[string]string{"Dockerfile":"harbur/test_web", "Dockerfile.backend": "harbur/test_backend"}
	assert.Equal(t, expected, c.GetImageNames(), "Should return image names")
}

func TestGetUnitTestCommands(t *testing.T) {
	options.config = "test/OneImage/captain.yml"
	c := NewConfig(options,false)
	expected := []string{"echo testing 1 web"}
	assert.Equal(t,expected, c.GetUnitTestCommands(), "Should return unit tests")
}

func TestFilterConfigEmpty(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options,false)
	expected := map[string]string{"Dockerfile":"harbur/test_web", "Dockerfile.backend": "harbur/test_backend"}
	assert.Equal(t, expected, c.GetImageNames(), "Should return complete list of image names")

	res := c.FilterConfig("")
	assert.True(t, res, "Should return true")
	assert.Equal(t, expected, c.GetImageNames(), "Should return complete list of image names")
}

func TestFilterConfigNonExistent(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options,false)
	expected := map[string]string{"Dockerfile":"harbur/test_web", "Dockerfile.backend": "harbur/test_backend"}
	assert.Equal(t, expected, c.GetImageNames(), "Should return complete list of image names")

	res := c.FilterConfig("nonexistent")
	assert.False(t, res, "Should return false")
	expected = map[string]string{}
	assert.Equal(t, expected, c.GetImageNames(), "Should return empty list")
}

func TestFilterConfigWeb(t *testing.T) {
	options.config = "test/Simple/captain.yml"
	c := NewConfig(options,false)
	expected := map[string]string{"Dockerfile":"harbur/test_web", "Dockerfile.backend": "harbur/test_backend"}
	assert.Equal(t, expected, c.GetImageNames(), "Should return complete list of image names")

	c.FilterConfig("web")
	expected = map[string]string{"Dockerfile":"harbur/test_web"}
	assert.Equal(t, expected, c.GetImageNames(), "Should return filtered list of image names")
}
