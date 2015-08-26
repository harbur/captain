package captain // import "github.com/harbur/captain"

import (
	"regexp"
	"strings"
)

func getRevision() string {
	res, _ := oneliner("git", "rev-parse", "--short", "HEAD")
	return res
}

func getBranches(all_branches bool) []string {
	// Labels (branches + tags)
	var labels =[]string{}

	branches_str, _ := oneliner("git", "name-rev", "--name-only", "HEAD")
	if (all_branches) {
		branches_str,_ = oneliner("git", "branch", "--no-column", "--contains", "HEAD")
	}

	var branches = make([]string, 5)
	if (branches_str != "") {
		// Remove asterisk from branches list
		r := regexp.MustCompile("[\\* ] ")
		branches_str = r.ReplaceAllString(branches_str, "")
		branches = strings.Split(branches_str, "\n")
		
		// Branches list is separated by spaces. Let's put it in an array
		labels=append(labels,branches...)
	}

	tags_str, _ := oneliner("git", "tag", "--points-at", "HEAD")

	if (tags_str != "") {
		tags := strings.Split(tags_str, "\n")
		debug("Active branches %s and tags %s", branches, tags)
		// Git tag list is separated by multi-lines. Let's put it in an array
		labels=append(labels,tags...)
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
