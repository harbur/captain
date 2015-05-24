package captain // import "github.com/harbur/captain/captain"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/v2/yaml"
)

type Config interface {
	GetImageNames() map[string]string
	GetUnitTestCommands() []string
}

type config struct {
	Build  build
	Test   map[string][]string
	Images []string
}

type build struct {
	Images map[string]string
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
	var res error
	if ext == ".json" {
		res = json.Unmarshal(data, &config)
	} else if ext == ".yml" || ext == ".yaml" {
		res = yaml.Unmarshal(data, &config)
	} else {
		panic(StatusError{errors.New("Unrecognized file extension"), 65})
	}
	if res != nil {
		res = displaySyntaxError(data, res)
		err("%s", res)
		os.Exit(InvalidCaptainYML)
	}
	return config
}

// NewConfig retus a new config based on given
// options.
// Containers will be ordered so that they can be
// brought up and down with Docker.
func NewConfig(options Options, forceOrder bool) Config {
	var conf *config
	for _, f := range configFiles(options) {
		if _, err := os.Stat(f); err == nil {
			conf = readConfig(f)
			break
		}
	}
	if conf == nil {
		info("No configuration found %v - inferring values", configFiles(options))
		conf = &config{}
		conf.Build.Images = make(map[string]string)

		conf.Build.Images = getDockerfiles()
	}

	var err error
	if err != nil {
		panic(StatusError{err, 78})
	}
	return conf
}

func (c *config) GetImageNames() map[string]string {
	return c.Build.Images
}

func (c *config) GetUnitTestCommands() []string {
	return c.Test["unit"]
}

// Global list, how can I pass it to the visitor pattern?
var imagesMap = make(map[string]string)

func getDockerfiles() map[string]string {
	filepath.Walk(".", visit)
	return imagesMap
}

func visit(path string, f os.FileInfo, err error) error {
	// Filename is "Dockerfile" or has "Dockerfile." prefix and is not a directory
	if (f.Name() == "Dockerfile" || strings.HasPrefix(f.Name(), "Dockerfile.")) && !f.IsDir() {
		// Get Parent Dirname
		absolute_path, _ := filepath.Abs(path)
		var image = strings.ToLower(filepath.Base(filepath.Dir(absolute_path)))
		imagesMap[path] = options.namespace + "/" + image + strings.ToLower(filepath.Ext(path))
		debug("Located %s will be used to create %s", path, imagesMap[path])
	}
	return nil
}
