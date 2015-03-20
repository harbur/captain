package captain // import "github.com/harbur/captain/captain"

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
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

				fmt.Printf("%s Building image %s\n", prefix("[CAPTAIN]"), info(image))

				execute("docker", "build", "-f", dockerfile, "-t", image, ".")

				var rev = getRevision()
				var imagename = image + ":" + rev
				fmt.Printf("%s Tagging image as %s\n", prefix("[CAPTAIN]"), info(imagename))
				execute("docker", "tag", "-f", image, imagename)

				var branch = getBranch()
				var branchname = image + ":" + branch
				fmt.Printf("%s Tagging image as %s\n", prefix("[CAPTAIN]"), info(branchname))
				execute("docker", "tag", "-f", image, branchname)
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
				fmt.Printf("%s Running unit test command: %s\n", prefix("[CAPTAIN]"), info(value))

				cmd := exec.Command("bash", "-c", value)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				cmd.Run()
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

	captainCmd.AddCommand(cmdBuild, cmdTest, cmdVersion)
	captainCmd.Execute()
}
