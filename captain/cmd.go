package captain // import "github.com/harbur/captain/captain"

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Options that are passed by CLI are mapped here for consumption
type Options struct {
	debug     bool
	force     bool
	namespace string
	config    string
	images    []string
}

var options Options

func handleCmd() {

	var cmdBuild = &cobra.Command{
		Use:   "build [image]",
		Short: "Builds the docker image(s) of your repository",
		Long:  `It will build the docker image(s) described on captain.yml in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			Build(config)
		},
	}

	var cmdTest = &cobra.Command{
		Use:   "test",
		Short: "Runs the tests",
		Long:  `It will execute the commands described on test section in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			// Build everything before testing
			Build(config)
			Test(config)
		},
	}

	var cmdPush = &cobra.Command{
		Use:   "push",
		Short: "Pushes the images to remote registry",
		Long:  `It will push the generated images to the remote registry.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			// Build everything before pushing
			Build(config)
			Push(config)
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  `Displays the version of Captain.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.2.0")
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

func getNamespace() string {
	return os.Getenv("USER")
}
