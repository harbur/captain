package captain // import "github.com/harbur/captain"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/v2/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config interface {
	GetImageNames() []string
	GetUnitTestCommands() []string
}

type config struct {
	Build  map[string][]string
	Test   map[string][]string
	Images []string
}

type Target []string

// configFiles returns a slice of
// files to read the config from.
// If the --config option was given,
// it will only use the given file.
func configFiles(options Options) []string {
	if len(options.config) > 0 {
		return []string{options.config}
	} else {
		return []string{"captain.json", "captain.yaml", "captain.yml"}
	}
}

// readConfig will read the config file
// and return the created config.
func readConfig(filename string) *config {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(StatusError{err, 74})
	}

	ext := filepath.Ext(filename)
	return unmarshal(data, ext)
}

// displaySyntaxError will display more information
// such as line and error type given an error and
// the data that was unmarshalled.
// Thanks to https://github.com/markpeek/packer/commit/5bf33a0e91b2318a40c42e9bf855dcc8dd4cdec5
func displaySyntaxError(data []byte, syntaxError error) (err error) {
	syntax, ok := syntaxError.(*json.SyntaxError)
	if !ok {
		err = syntaxError
		return
	}
	newline := []byte{'\x0a'}
	space := []byte{' '}

	start, end := bytes.LastIndex(data[:syntax.Offset], newline)+1, len(data)
	if idx := bytes.Index(data[start:], newline); idx >= 0 {
		end = start + idx
	}

	line, pos := bytes.Count(data[:start], newline)+1, int(syntax.Offset)-start-1

	err = fmt.Errorf("\nError in line %d: %s \n%s\n%s^", line, syntaxError, data[start:end], bytes.Repeat(space, pos))
	return
}

// unmarshal converts either JSON
// or YAML into a config object.
func unmarshal(data []byte, ext string) *config {
	var config *config
	var err error
	if ext == ".json" {
		err = json.Unmarshal(data, &config)
	} else if ext == ".yml" || ext == ".yaml" {
		err = yaml.Unmarshal(data, &config)
	} else {
		panic(StatusError{errors.New("Unrecognized file extension"), 65})
	}
	if err != nil {
		err = displaySyntaxError(data, err)
		panic(StatusError{err, 65})
	}
	return config
}

// NewConfig retus a new config based on given
// options.
// Containers will be ordered so that they can be
// brought up and down with Docker.
func NewConfig(options Options, forceOrder bool) Config {
	var config *config
	for _, f := range configFiles(options) {
		if _, err := os.Stat(f); err == nil {
			config = readConfig(f)
			break
		}
	}
	if config == nil {
		panic(StatusError{fmt.Errorf("No configuration found %v", configFiles(options)), 78})
	}

	var err error
	if err != nil {
		panic(StatusError{err, 78})
	}
	return config
}

func (c *config) GetImageNames() []string {
	// fmt.Printf("%#v\n", c.Build["images"])
	return c.Build["images"]
}

func (c *config) GetUnitTestCommands() []string {
	fmt.Printf("%#v\n", c.Test["unit"])
	return c.Test["unit"]
}
