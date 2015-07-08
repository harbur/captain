package captain // import "github.com/harbur/captain/captain"

import (
    "strings"
    "regexp"
)

func getRevision() string {
	res,_ := oneliner("git", "rev-parse", "--short", "HEAD")
	return res
}

func getBranch() string {
	branch,_ := oneliner("git", "name-rev", "--name-only", "HEAD")
	tag,err := oneliner("git", "tag", "--points-at", "HEAD")
	if (err ==nil && tag != "") {
		branch = tag
	}

	// Remove start of "heads/origin" if exist
	r := regexp.MustCompile("^heads\\/origin\\/")
	branch = r.ReplaceAllString(branch, "")

	// Remove start of "remotes/origin" if exist
	r = regexp.MustCompile("^remotes\\/origin\\/")
	branch = r.ReplaceAllString(branch, "")

	// Replace all "/" with "."
	branch = strings.Replace(branch, "/", ".", -1)

	// Replace all "~" with "."
	branch = strings.Replace(branch, "~", ".", -1)

	return branch
}

func isDirty() bool {
	res,_ := oneliner("git", "status", "--porcelain")
	return len(res) > 0
}

func isGit() bool {
	res := getRevision()
	return len(res) > 0
}
