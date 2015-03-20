package captain // import "github.com/harbur/captain/captain"

import (
	"bytes"
	"os/exec"
	"strings"
)

func getRevision() string {
	var buff bytes.Buffer
	gitCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	gitCmd.Stdout = &buff
	gitCmd.Run()
	var rev = strings.TrimSpace(buff.String())

	return rev
}
