package captain // import "github.com/harbur/captain"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config represents the information stored at captain.yml. It keeps information about images and unit tests.
type Config interface {
	FilterConfig(filter string) bool
	GetApp(app string) App
	GetApps() []App
}

type configV1 struct {
	Build  build
	Test   map[string][]string
	Images []string
	Root   []string
}

type build struct {
	Images map[string]string
}

type config map[string]App

var configOrder *yaml.MapSlice

// App struct
type App struct {
	Build     string
	Image     string
	Pre       []string
	Post      []string
	Test      []string
	Build_arg map[string]string
}

// configFile returns the file to read the config from.
// If the --config option was given,
// it will only use the given file.
func configFile(path string) string {
	if len(path) > 0 {
		return path
	}
	return "captain.yml"
}

// readConfig will read the config file
// and return the created config.
func readConfig(filename string) *config {
	data, err := ioutil.ReadFile(filename)
	os.Chdir(filepath.Dir(filename))
	if err != nil {
		panic(StatusError{err, 74})
	}
	return unmarshal(data)
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
func unmarshal(data []byte) *config {
	var configV1 *configV1
	res := yaml.Unmarshal(data, &configV1)
	if len(configV1.Build.Images) > 0 {
		err("Old %s format detected! Please check the https://github.com/harbur/captain how to upgrade", "captain.yml")
		os.Exit(-1)
	}

	var config *config
	res = yaml.Unmarshal(data, &config)

	if res != nil {
		res = displaySyntaxError(data, res)
		err("%s", res)
		os.Exit(InvalidCaptainYML)
	}

	// We re-import it as MapSlice to keep order of apps
	res = yaml.Unmarshal(data, &configOrder)

	if res != nil {
		res = displaySyntaxError(data, res)
		err("%s", res)
		os.Exit(InvalidCaptainYML)
	}

	return config
}

// NewConfig returns a new Config instance based on the reading the captain.yml
// file at path.
// Containers will be ordered so that they can be
// brought up and down with Docker.
func NewConfig(namespace, path string, forceOrder bool) Config {
	var conf *config
	f := configFile(path)
	if _, err := os.Stat(f); err == nil {
		conf = readConfig(f)
	}

	if conf == nil {
		info("No configuration found %v - inferring values", configFile(path))
		autoconf := make(config)
		conf = &autoconf
		dockerfiles := getDockerfiles(namespace)
		for build, image := range dockerfiles {
			autoconf[image] = App{Build: build, Image: image}
		}
	}

	var err error
	if err != nil {
		panic(StatusError{err, 78})
	}
	return conf
}

// GetApps returns a list of Apps
func (c *config) GetApps() []App {
	var cc = *c
	var apps []App
	if configOrder != nil {
		for _, v := range *configOrder {
			if val, ok := cc[v.Key.(string)]; ok {
				apps = append(apps, val)
			}
		}
	} else {
		for _, v := range *c {
			apps = append(apps, v)
		}
	}

	return apps
}

func (c *config) FilterConfig(filter string) bool {
	if filter != "" {
		res := false
		for key := range *c {
			if key == filter {
				res = true
			} else {
				delete(*c, key)
			}
		}
		return res
	}
	return true
}

// GetApp returns App configuration
func (c *config) GetApp(app string) App {
	for key, k := range *c {
		if key == app {
			return k
		}
	}
	return App{}
}

// Global list, how can I pass it to the visitor pattern?
var imagesMap = make(map[string]string)

func getDockerfiles(namespace string) map[string]string {
	filepath.Walk(".", visit(namespace))
	return imagesMap
}

func visit(namespace string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		// Filename is "Dockerfile" or has "Dockerfile." prefix and is not a directory
		if (f.Name() == "Dockerfile" || strings.HasPrefix(f.Name(), "Dockerfile.")) && !f.IsDir() {
			// Get Parent Dirname
			absolutePath, _ := filepath.Abs(path)
			var image = strings.ToLower(filepath.Base(filepath.Dir(absolutePath)))
			delimitedNamespace := namespace
			if !strings.Contains(namespace, "/") {
				delimitedNamespace += "/"
			}
			imagesMap[path] = delimitedNamespace + image + strings.ToLower(filepath.Ext(path))
			debug("Located %s will be used to create %s", path, imagesMap[path])
		}
		return nil
	}
}
