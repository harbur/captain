package captain // import "github.com/harbur/captain"

const (
	// BuildFailed represents a build failure
	BuildFailed = 1

	// TagFailed represents a failure to tag a docker image
	TagFailed = 2

	// NonExistImage represents the existance of a docker image tag
	NonExistImage = 3

	// TestFailed represents test failure
	TestFailed = 5

	// NoGit represents lack of a git repository
	NoGit = 6

	// GitDirty represents existence of local git changes
	GitDirty = 7

	// InvalidCaptainYML represents an invalid captain.yml format
	InvalidCaptainYML = 8

	// NoDockerfiles represents lack of Dockerfile(s) on current and subdirectories.
	NoDockerfiles = 9

	// OldFormat represents old format of captain.yml
	OldFormat = 10
)
