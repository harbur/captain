package captain // import "github.com/harbur/captain"

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harbur/captain/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

const defaultEndpoint = "unix:///var/run/docker.sock"

var client *docker.Client

func init() {
	var err error
	client, err = docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
}

func buildImage(app App, tag string, force bool) error {
	info("Building image %s:%s", app.Image, tag)

	// Nasty issue with CircleCI https://github.com/docker/docker/issues/4897
	if os.Getenv("CIRCLECI") == "true" {
		info("Running at %s environment...", "CIRCLECI")
		execute("docker", "build", "-t", app.Image+":"+tag, filepath.Dir(app.Build))
		return nil
	}

	opts := docker.BuildImageOptions{
		Name:                app.Image + ":" + tag,
		Dockerfile:          filepath.Base(app.Build),
		NoCache:             force,
		SuppressOutput:      false,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		OutputStream:        os.Stdout,
		ContextDir:          filepath.Dir(app.Build),
	}
	err := client.BuildImage(opts)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return err
}

func tagImage(app App, origin string, tag string) error {
	if tag != "" {
		info("Tagging image %s:%s as %s:%s", app.Image, origin, app.Image, tag)
		opts := docker.TagImageOptions{Repo: app.Image, Tag: tag, Force: true}
		err := client.TagImage(app.Image+":"+origin, opts)
		if err != nil {
			fmt.Printf("%s", err)
		}
		return err
	}

	debug("Skipping tag of %s - no git repository", app.Image)

	return nil
}

func removeImage(name string) error {
	return client.RemoveImage(name)
}

/**
 * Retrieves a list of existing Images for the specific App.
 */
func getImages(app App) []docker.APIImages {
	debug("Getting images %s", app.Image)
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false, Filter: app.Image})
	return imgs
}

func imageExist(app App, tag string) bool {
	repo := app.Image + ":" + tag
	image, _ := client.InspectImage(repo)
	if image != nil {
		return true
	}
	return false
}
