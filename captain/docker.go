package captain // import "github.com/harbur/captain/captain"

import (
	"os"

	"github.com/fsouza/go-dockerclient"
)

var endpoint = "unix:///var/run/docker.sock"
var client, _ = docker.NewClient(endpoint)

func buildImage(dockerfile string, image string, tag string) {
	info("Building image %s:%s", image, tag)

	opts := docker.BuildImageOptions{
		Name:                image + ":" + tag,
		NoCache:             false,
		SuppressOutput:      false,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		OutputStream:        os.Stdout,
		ContextDir:          ".",
	}
	client.BuildImage(opts)
}

func tagImage(repo string, origin string, tag string) {
	info("Tagging image %s:%s as %s:%s", repo, origin, repo, tag)
	// var imageID = getImageID(repo, origin)
	opts := docker.TagImageOptions{Repo: repo, Tag: tag, Force: true}
	client.TagImage(repo, opts)
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
