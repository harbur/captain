package captain // import "github.com/harbur/captain/captain"

import (
    "strings"
    "regexp"
)

func getRevision() string {
	return oneliner("git", "rev-parse", "--short", "HEAD")
}

func getBranch() string {
	branch := oneliner("git", "rev-parse", "--abbrev-ref", "HEAD")

	// Remove start of "heads/origin" if exist
	r := regexp.MustCompile("^heads\\/origin\\/")
	branch = r.ReplaceAllString(branch, "")

	// Replace all "/" with "."
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
