package captain // import "github.com/harbur/captain/captain"

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Options struct {
	debug  bool
	config string
	images []string
}

var options Options

var (
	BUILD_FAILED = 1
	TAG_FAILED   = 2
)

func handleCmd() {

	var cmdBuild = &cobra.Command{
		Use:   "build [image]",
		Short: "Builds the docker image(s) of your repository",
		Long:  `It will build the docker image(s) described on captain.yml in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			var images = config.GetImageNames()

			if len(args) == 1 {
				images = filterImages(images, args[0])
			}
			var rev = getRevision()

			for _, value := range images {
				s := strings.Split(value, "=")
				dockerfile, image := s[0], s[1]

				// Skip build if there are no local changes and the commit is already built
				if !isDirty() && imageExist(image, rev) {
					info("Skipping build of %s:%s - image is already built", image, rev)

					// Tag commit image
					tagImage(image, rev, "latest")

					// Tag branch image
					var branch = getBranch()
					if branch == "HEAD" {
						debug("Skipping tag of %s in detached mode", image)
					} else {
						tagImage(image, rev, branch)
					}

				} else {
					// Build latest image
					err := buildImage(dockerfile, image, "latest")
					if err != nil {
						os.Exit(BUILD_FAILED)
					}
					if isDirty() {
						debug("Skipping tag of %s:%s - local changes exist", image, rev)
					} else {
						// Tag commit image
						tagImage(image, "latest", rev)

						// Tag branch image
						var branch = getBranch()
						if branch == "HEAD" {
							debug("Skipping tag of %s in detached mode", image)
						} else {
							err := tagImage(image, "latest", branch)
							if err != nil {
								os.Exit(TAG_FAILED)
							}
						}
					}

				}

			}
		},
	}

	var cmdTest = &cobra.Command{
		Use:   "test",
		Short: "Runs the unit tests",
		Long:  `It will execute the commands described on unit testing in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			for _, value := range config.GetUnitTestCommands() {
				info("Running unit test command: %s", value)
				execute("bash", "-c", value)
			}
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  `Displays the version of Crane.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.0.1")
		},
	}

	var captainCmd = &cobra.Command{
		Use:   "captain",
		Short: "captain - build tool for Docker focused on CI/CD",
		Long: `
Captain, the CLI build tool for Docker made for Continuous Integration / Continuous Delivery.

It works by reading captain.yaml file which describes how to build, test, push and release the docker image(s) of your repository.`,
	}

	captainCmd.PersistentFlags().BoolVarP(&options.debug, "debug", "D", false, "Enable debug mode")
	captainCmd.AddCommand(cmdBuild, cmdTest, cmdVersion)
	captainCmd.Execute()
}

func filterImages(images []string, arg string) []string {
	for _, value := range images {
		s := strings.Split(value, "=")
		_, image := s[0], s[1]
		if image == arg {
			return []string{value}
		}
	}
	err("Build image %s is not defined", arg)
	os.Exit(-1)
	return []string{}
}
