package captain // import "github.com/harbur/captain/captain"

import (
	"os"
	"os/exec"
)

func execute(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
	return cmd
}
