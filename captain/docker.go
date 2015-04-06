package captain // import "github.com/harbur/captain/captain"

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

var endpoint = "unix:///var/run/docker.sock"
var client, _ = docker.NewClient(endpoint)

func buildImage(dockerfile string, image string, tag string) error {
	info("Building image %s:%s", image, tag)

	// Nasty issue with CircleCI https://github.com/docker/docker/issues/4897
	if os.Getenv("CIRCLECI") == "true" {
		info("Running at %s environment...", "CIRCLECI")
		execute("docker", "build", "-t", image+":"+tag, ".")
		return nil
	} else {
		opts := docker.BuildImageOptions{
			Name:                image + ":" + tag,
			Dockerfile:          dockerfile,
			NoCache:             options.force,
			SuppressOutput:      false,
			RmTmpContainer:      true,
			ForceRmTmpContainer: true,
			OutputStream:        os.Stdout,
			ContextDir:          ".",
		}
		err := client.BuildImage(opts)
		if err != nil {
			fmt.Printf("%s", err)
		}
		return err
	}
}

func tagImage(repo string, origin string, tag string) error {
	if tag != "" {
		info("Tagging image %s:%s as %s:%s", repo, origin, repo, tag)
		// var imageID = getImageID(repo, origin)
		opts := docker.TagImageOptions{Repo: repo, Tag: tag, Force: true}
		err := client.TagImage(repo, opts)
		if err != nil {
			fmt.Printf("%s", err)
		}
		return err
	} else {
		debug("Skipping tag of %s - no git repository", repo)
	}
	return nil
}

func getImageID(repo string, tag string) string {
	images, _ := client.ListImages(docker.ListImagesOptions{})
	for _, image := range images {
		for _, b := range image.RepoTags {
			if b == repo+":"+tag {
				return image.ID
			}
		}
	}
	return ""
}

func imageExist(repo string, tag string) bool {
	images, _ := client.ListImages(docker.ListImagesOptions{})
	for _, image := range images {
		for _, b := range image.RepoTags {
			if b == repo+":"+tag {
				return true
			}
		}
	}
	return false
}
