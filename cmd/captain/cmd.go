package main

import (
	"fmt"
	"os"

	"github.com/harbur/captain"
	"github.com/harbur/captain/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/harbur/captain/Godeps/_workspace/src/github.com/spf13/cobra"
)

// Options that are passed by CLI are mapped here for consumption
type Options struct {
	debug        bool
	force        bool
	all_branches bool
	namespace    string
	config       string
	images       []string
}

var options Options

func handleCmd() {

	var cmdBuild = &cobra.Command{
		Use:   "build [image]",
		Short: "Builds the docker image(s) of your repository",
		Long:  `It will build the docker image(s) described on captain.yml in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := captain.NewConfig(options.namespace, options.config, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			buildOpts := captain.BuildOptions{
				Config: config,
				Force:  options.force,
				All_branches:  options.all_branches,
			}

			captain.Build(buildOpts)
		},
	}

	var cmdTest = &cobra.Command{
		Use:   "test",
		Short: "Runs the tests",
		Long:  `It will execute the commands described on test section in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := captain.NewConfig(options.namespace, options.config, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			buildOpts := captain.BuildOptions{
				Config: config,
				Force:  options.force,
				All_branches:  options.all_branches,
			}

			// Build everything before testing
			captain.Build(buildOpts)
			captain.Test(buildOpts)
		},
	}

	var cmdPush = &cobra.Command{
		Use:   "push",
		Short: "Pushes the images to remote registry",
		Long:  `It will push the generated images to the remote registry.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := captain.NewConfig(options.namespace, options.config, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			buildOpts := captain.BuildOptions{
				Config: config,
				Force:  options.force,
				All_branches:  options.all_branches,
			}

			// Build everything before pushing
			captain.Build(buildOpts)
			captain.Push(buildOpts)
		},
	}

	var cmdPull = &cobra.Command{
		Use:   "pull",
		Short: "Pulls the images from remote registry",
		Long:  `It will pull the images from the remote registry.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := captain.NewConfig(options.namespace, options.config, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			buildOpts := captain.BuildOptions{
				Config: config,
				Force:  options.force,
				All_branches:  options.all_branches,
			}

			captain.Pull(buildOpts)
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  `Displays the version of Captain.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.6.0")
		},
	}

	var captainCmd = &cobra.Command{
		Use:   "captain",
		Short: "captain - build tool for Docker focused on CI/CD",
		Long: `
Captain, the CLI build tool for Docker made for Continuous Integration / Continuous Delivery.

It works by reading captain.yaml file which describes how to build, test, push and release the docker image(s) of your repository.`,
	}

	captainCmd.PersistentFlags().BoolVarP(&captain.Debug, "debug", "D", false, "Enable debug mode")
	captainCmd.PersistentFlags().StringVarP(&options.namespace, "namespace", "N", getNamespace(), "Set default image namespace")
	captainCmd.PersistentFlags().BoolVarP(&color.NoColor, "no-color", "n", false, "Disable color output")
	captainCmd.PersistentFlags().BoolVarP(&options.all_branches, "all-branches", "B", false, "Build all branches on specific commit instead of just working branch")
	cmdBuild.Flags().BoolVarP(&options.force, "force", "f", false, "Force build even if image is already built")
	captainCmd.AddCommand(cmdBuild, cmdTest, cmdPush, cmdPull, cmdVersion)
	captainCmd.Execute()
}

func getNamespace() string {
	return os.Getenv("USER")
}
