package captain // import "github.com/harbur/captain/captain"

import (
    "strings"
)

func getRevision() string {
	return oneliner("git", "rev-parse", "--short", "HEAD")
}

func getBranch() string {
	branch := oneliner("git", "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.Replace(branch, "/", ".", -1)
	return branch
}

func isDirty() bool {
	var res = oneliner("git", "status", "--porcelain")
	return len(res) > 0
}

func isGit() bool {
	var res = getRevision()
	return len(res) > 0
}
