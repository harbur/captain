package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/harbur/captain"
	"github.com/spf13/cobra"
)

// Options that are passed by CLI are mapped here for consumption
type Options struct {
	debug     bool
	force     bool
	long_sha  bool
	namespace string
	config    string
	images    []string
	tag       string

	// Options to define the docker tags context
	all_branches bool
	branch_tags  bool
	commit_tags  bool
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
				Config:       config,
				Tag:          options.tag,
				Force:        options.force,
				All_branches: options.all_branches,
				Long_sha:     options.long_sha,
				Branch_tags:  options.branch_tags,
				Commit_tags:  options.commit_tags,
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
				Config:       config,
				Tag:          options.tag,
				Force:        options.force,
				All_branches: options.all_branches,
				Long_sha:     options.long_sha,
				Branch_tags:  options.branch_tags,
				Commit_tags:  options.commit_tags,
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
				Config:       config,
				Tag:          options.tag,
				Force:        options.force,
				All_branches: options.all_branches,
				Long_sha:     options.long_sha,
				Branch_tags:  options.branch_tags,
				Commit_tags:  options.commit_tags,
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
				Config:       config,
				Tag:          options.tag,
				Force:        options.force,
				All_branches: options.all_branches,
				Long_sha:     options.long_sha,
				Branch_tags:  options.branch_tags,
				Commit_tags:  options.commit_tags,
			}

			captain.Pull(buildOpts)
		},
	}

	var cmdPurge = &cobra.Command{
		Use:   "purge",
		Short: "Purges the stale images",
		Long:  `It will purge the stale images. Stale image is an image that is not the latest of at least one branch.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := captain.NewConfig(options.namespace, options.config, true)

			if len(args) == 1 {
				config.FilterConfig(args[0])
			}

			buildOpts := captain.BuildOptions{
				Config:       config,
				Force:        options.force,
				All_branches: options.all_branches,
				Long_sha:     options.long_sha,
			}

			captain.Purge(buildOpts)
		},
	}

	var cmdSelfUpdate = &cobra.Command{
		Use:   "self-update",
		Short: "Updates Captain to the last version",
		Long:  `Updates Captain to the last available version.`,
		Run: func(cmd *cobra.Command, args []string) {
			captain.SelfUpdate()
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  `Displays the version of Captain.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v1.1.0")
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
	captainCmd.PersistentFlags().BoolVarP(&options.long_sha, "long-sha", "l", false, "Use the long git commit SHA when referencing revisions")

	cmdBuild.Flags().BoolVarP(&options.force, "force", "f", false, "Force build even if image is already built")
	cmdBuild.Flags().BoolVarP(&options.all_branches, "all-branches", "B", false, "Build all branches on specific commit instead of just working branch")
	cmdBuild.Flags().StringVarP(&options.tag, "tag", "t", "", "Tag version")

	cmdPull.Flags().BoolVarP(&options.all_branches, "all-branches", "B", false, "Pull all branches on specific commit instead of just working branch")
	cmdPull.Flags().BoolVarP(&options.branch_tags, "branch-tags", "b", true, "Pull the 'branch' docker tags")
	cmdPull.Flags().BoolVarP(&options.commit_tags, "commit-tags", "c", false, "Pull the 'commit' docker tags")
	cmdPull.Flags().StringVarP(&options.tag, "tag", "t", "", "Tag version")

	cmdPush.Flags().BoolVarP(&options.all_branches, "all-branches", "B", false, "Push all branches on specific commit instead of just working branch")
	cmdPush.Flags().BoolVarP(&options.branch_tags, "branch-tags", "b", true, "Push the 'branch' docker tags")
	cmdPush.Flags().BoolVarP(&options.commit_tags, "commit-tags", "c", false, "Push the 'commit' docker tags")
	cmdPush.Flags().StringVarP(&options.tag, "tag", "t", "", "Tag version")

	cmdPurge.Flags().BoolVarP(&options.force, "dangling", "d", false, "Remove dangling images")

	captainCmd.AddCommand(cmdBuild, cmdTest, cmdPush, cmdPull, cmdVersion, cmdPurge, cmdSelfUpdate)
	captainCmd.Execute()
}

func getNamespace() string {
	return os.Getenv("USER")
}
