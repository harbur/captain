package captain // import "github.com/harbur/captain/captain"

func getRevision() string {
	return oneliner("git", "rev-parse", "--short", "HEAD")
}

func getBranch() string {
	return oneliner("git", "rev-parse", "--abbrev-ref", "HEAD")
}
