package captain // import "github.com/harbur/captain/captain"

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Options struct {
	debug  bool
	force  bool
	config string
	images []string
}

var options Options

var (
	BUILD_FAILED   = 1
	TAG_FAILED     = 2
	NONEXIST_IMAGE = 3
	NO_CAPTAIN_YML = 4
	TEST_FAILED    = 5
)

func handleCmd() {

	var cmdBuild = &cobra.Command{
		Use:   "build [image]",
		Short: "Builds the docker image(s) of your repository",
		Long:  `It will build the docker image(s) described on captain.yml in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			var images string
			if len(args) == 1 {
				images = args[0]
			}

			Build(config, images)
		},
	}

	var cmdTest = &cobra.Command{
		Use:   "test",
		Short: "Runs the unit tests",
		Long:  `It will execute the commands described on unit testing in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			// Build everything before testing
			Build(config, "")
			Test(config, "")
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
	cmdBuild.Flags().BoolVarP(&options.force, "force", "f", false, "Force build even if image is already built")
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
	os.Exit(NONEXIST_IMAGE)
	return []string{}
}
