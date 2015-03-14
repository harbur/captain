package captain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Options struct {
	config string
	images []string
}

var options Options

func handleCmd() {

	var cmdBuild = &cobra.Command{
		Use:   "build",
		Short: "Builds the docker image(s) of your repository",
		Long:  `It will build the docker image(s) described on captain.yml in order they appear on file.`,
		Run: func(cmd *cobra.Command, args []string) {
			config := NewConfig(options, true)

			for _, value := range config.GetImageNames() {
				s := strings.Split(value, "=")
				dockerfile, image := s[0], s[1]
				fmt.Println(dockerfile + image)

				var buff bytes.Buffer

				gitCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
				gitCmd.Stdout = &buff
				gitCmd.Run()
				fmt.Println(buff.String())
				var rev = strings.TrimSpace(buff.String())

				fmt.Println("Building image " + image)
				cmd := exec.Command("docker", "build", "-f", dockerfile, "-t", image+":"+rev, ".")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				cmd.Run()
				// cmd := gitCmd
				if !cmd.ProcessState.Success() {
					status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
					panic(StatusError{errors.New(cmd.ProcessState.String()), status})
				}
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

	captainCmd.AddCommand(cmdBuild, cmdVersion)
	captainCmd.Execute()
}
