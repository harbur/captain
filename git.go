package captain // import "github.com/harbur/captain"

import (
	"regexp"
	"strings"
)

func getRevision(long_sha bool) string {
	params := []string{"rev-parse"}
	if !long_sha {
		params = append(params, "--short")
	}

	params = append(params, "HEAD")
	res, _ := oneliner("git", params...)
	return res
}

func getBranches(all_branches bool) []string {
	// Labels (branches + tags)
	var labels = []string{}

	branches_str, _ := oneliner("git", "name-rev", "--name-only", "--exclude=tags/*", "HEAD")
	if all_branches {
		branches_str, _ = oneliner("git", "branch", "--no-column", "--points-at", "HEAD")
	}

	var branches = make([]string, 5)
	if branches_str != "" {
		// Remove asterisk from branches list
		r := regexp.MustCompile("[\\* ] ")
		branches_str = r.ReplaceAllString(branches_str, "")
		// Remove "(HEAD detached at..." if not on a branch
		r = regexp.MustCompile("\\(HEAD detached at .*\\)")
		branches_str = r.ReplaceAllString(branches_str, "")
		branches = strings.Split(branches_str, "\n")

		// Branches list is separated by spaces. Let's put it in an array
		labels = append(labels, branches...)
	}

	tags_str, _ := oneliner("git", "tag", "--points-at", "HEAD")

	if tags_str != "" {
		tags := strings.Split(tags_str, "\n")
		debug("Active branches %s and tags %s", branches, tags)
		// Git tag list is separated by multi-lines. Let's put it in an array
		labels = append(labels, tags...)
	}

	var cleanLabels = []string{}
	for key := range labels {
		var label = labels[key]
		// Remove start of "heads/origin" if exist
		label = regexp.MustCompile("^heads\\/origin\\/").ReplaceAllString(label, "")

		// Remove start of "remotes/origin" if exist
		label = regexp.MustCompile("^remotes\\/origin\\/").ReplaceAllString(label, "")

		// Replace all "/" with "."
		label = strings.Replace(label, "/", ".", -1)

		// Replace all "~" with "."
		label = strings.Replace(label, "~", ".", -1)

		if label != "" {
			cleanLabels = append(cleanLabels, label)
		}
	}
	return cleanLabels
}

func isDirty() bool {
	res, _ := oneliner("git", "status", "--porcelain")
	return len(res) > 0
}

func isGit() bool {
	res := getRevision(false)
	return len(res) > 0
}
