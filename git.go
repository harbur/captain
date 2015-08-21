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
	// Labels (branches + tags)
	var labels =[]string{}

	branches,_ := oneliner("git", "branch", "--no-column", "--contains", "HEAD")
	if (branches != "") {
		// Remove asterisk from branches list
		r := regexp.MustCompile("[\\* ] ")
		branches = r.ReplaceAllString(branches, "")
		// Branches list is separated by spaces. Let's put it in an array
		labels=append(labels,strings.Split(branches, "\n")...)
	}

	tags, _ := oneliner("git", "tag", "--points-at", "HEAD")

	if (tags != "") {
		// Git tag list is separated by multi-lines. Let's put it in an array
		labels=append(labels,strings.Split(tags, "\n")...)
	}

	for key := range labels {
		// Remove start of "heads/origin" if exist
		r := regexp.MustCompile("^heads\\/origin\\/")
		labels[key] = r.ReplaceAllString(labels[key], "")

		// Remove start of "remotes/origin" if exist
		r = regexp.MustCompile("^remotes\\/origin\\/")
		labels[key] = r.ReplaceAllString(labels[key], "")

		// Replace all "/" with "."
		labels[key] = strings.Replace(labels[key], "/", ".", -1)

		// Replace all "~" with "."
		labels[key] = strings.Replace(labels[key], "~", ".", -1)
	}

	return labels
}

func isDirty() bool {
	res, _ := oneliner("git", "status", "--porcelain")
	return len(res) > 0
}

func isGit() bool {
	res := getRevision()
	return len(res) > 0
}
