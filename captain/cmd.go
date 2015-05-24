package captain // import "github.com/harbur/captain/captain"

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	debug     bool
	force     bool
	namespace string
	config    string
	images    []string
}

var options Options

var (
	BuildFailed       = 1
	TagFailed         = 2
	NonexistImage     = 3
	NoCaptainYML      = 4
	TestFailed        = 5
	NoGit             = 6
	GitDirty          = 7
	InvalidCaptainYML = 8
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

	var cmdPush = &cobra.Command{
		Use:   "push",
		Short: "Pushes the images to remote registry",
		Long:  `It will push the generated images to the remote registry.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			// Build everything before pushing
			Build(config, "")
			Push(config, "")
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
	captainCmd.PersistentFlags().StringVarP(&options.namespace, "namespace", "N", getNamespace(), "Set default image namespace")
	cmdBuild.Flags().BoolVarP(&options.force, "force", "f", false, "Force build even if image is already built")
	captainCmd.AddCommand(cmdBuild, cmdTest, cmdPush, cmdVersion)
	captainCmd.Execute()
}

func filterImages(images map[string]string, arg string) map[string]string {
	for key, image := range images {
		if image == arg {
			return map[string]string{key: image}
		}
	}
	err("Build image %s is not defined", arg)
	os.Exit(NonexistImage)
	return map[string]string{}
}

func getNamespace() string {
	return os.Getenv("USER")
}
