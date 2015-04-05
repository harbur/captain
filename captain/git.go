package captain // import "github.com/harbur/captain/captain"

func getRevision() string {
	return oneliner("git", "rev-parse", "--short", "HEAD")
}

func getBranch() string {
	return oneliner("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func isDirty() bool {
	var res = oneliner("git", "status", "--porcelain")
	return len(res) > 0
}

func isGit() bool {
	var res = getRevision()
	return len(res) > 0
}
