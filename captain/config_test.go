package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFiles(t *testing.T) {
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
	options.config = "../captain.yml"
	c := NewConfig(options,false)
	expected := []string{"docker run harbur/captain-test go test github.com/harbur/captain/captain"}
	assert.Equal(t,expected, c.GetUnitTestCommands(), "Should return unit tests")
}