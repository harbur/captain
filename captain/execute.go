package captain // import "github.com/harbur/captain/captain"

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func execute(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
	return cmd
}

func oneliner(name string, arg ...string) string {
	var buff bytes.Buffer
	gitCmd := exec.Command(name, arg...)
	gitCmd.Stdout = &buff
	gitCmd.Run()
	return strings.TrimSpace(buff.String())
}
