package captain // import "github.com/harbur/captain"

import (
	"regexp"
	"strings"
)

func getRevision() string {
	res, _ := oneliner("git", "rev-parse", "--short", "HEAD")
	return res
}

func getBranches() []string {
	branch, _ := oneliner("git", "rev-parse", "--abbrev-ref", "HEAD")
	tag, err := oneliner("git", "tag", "--points-at", "HEAD")
	if err == nil && tag != "" {
		branch = tag
	}

	// Git tag list is separated in multi-lines. Let's put it in an array
	branches := strings.Split(branch, "\n")

	for key := range branches {
		// Remove start of "heads/origin" if exist
		r := regexp.MustCompile("^heads\\/origin\\/")
		branches[key] = r.ReplaceAllString(branches[key], "")

		// Remove start of "remotes/origin" if exist
		r = regexp.MustCompile("^remotes\\/origin\\/")
		branches[key] = r.ReplaceAllString(branches[key], "")

		// Replace all "/" with "."
		branches[key] = strings.Replace(branches[key], "/", ".", -1)

		// Replace all "~" with "."
		branches[key] = strings.Replace(branches[key], "~", ".", -1)
	}

	return branches
}

func isDirty() bool {
	res, _ := oneliner("git", "status", "--porcelain")
	return len(res) > 0
}

func isGit() bool {
	res := getRevision()
	return len(res) > 0
}
