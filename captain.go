package captain // import "github.com/harbur/captain"

import (
	"os"
)

// Debug can be turned on to enable debug mode.
var Debug bool

// StatusError provides error code and id
type StatusError struct {
	error  error
	status int
}

// Pre function executes commands on pre section before build
func Pre(config Config, app App) {
	for _, value := range app.Pre {
		info("Running pre command: %s", value)
		res := execute("bash", "-c", value)
		if res != nil {
			err("Pre execution returned non-zero status")
			os.Exit(TestFailed)
		}
	}
}

// Post function executes commands on pre section after build
func Post(config Config, app App) {
	for _, value := range app.Post {
		info("Running post command: %s", value)
		res := execute("bash", "-c", value)
		if res != nil {
			err("Post execution returned non-zero status")
			os.Exit(TestFailed)
		}
	}
}

type BuildOptions struct {
	Config Config
	Force  bool
	All_branches bool
}

// Build function compiles the Containers of the project
func Build(opts BuildOptions) {
	config := opts.Config

	var rev = getRevision()

	// For each App
	for _, app := range config.GetApps() {
		// If no Git repo exist
		if !isGit() {
			// Perfoming [build latest]
			debug("No local git repository found, just building latest")

			// Execute Pre commands
			Pre(config, app)

			// Build latest image
			res := buildImage(app, "latest", opts.Force)
			if res != nil {
				os.Exit(BuildFailed)
			}
		} else {
			// Skip build if there are no local changes and the commit is already built
			if !isDirty() && imageExist(app, rev) && !opts.Force {
				// Performing [skip rev|tag rev@latest|tag rev@branch]
				info("Skipping build of %s:%s - image is already built", app.Image, rev)

				// Tag commit image
				tagImage(app, rev, "latest")

				// Tag branch image
				for _,branch := range getBranches(opts.All_branches) {
					res := tagImage(app, rev, branch)
					if res != nil {
						os.Exit(TagFailed)
					}
				}

			} else {
				// Performing [build latest|tag latest@rev|tag latest@branch]

				// Execute Pre commands
				Pre(config, app)

				// Build latest image
				res := buildImage(app, "latest", opts.Force)
				if res != nil {
					os.Exit(BuildFailed)
				}
				if isDirty() {
					debug("Skipping tag of %s:%s - local changes exist", app.Image, rev)
				} else {
					// Tag commit image
					tagImage(app, "latest", rev)

					// Tag branch image
					for _,branch := range getBranches(opts.All_branches) {
						res := tagImage(app, "latest", branch)
						if res != nil {
							os.Exit(TagFailed)
						}
					}
				}
			}
		}
		Post(config, app)
	}
}

// Test function executes the tests of the project
func Test(opts BuildOptions) {
	config := opts.Config

	for _, app := range config.GetApps() {
		for _, value := range app.Test {
			info("Running test command: %s", value)
			res := execute("bash", "-c", value)
			if res != nil {
				err("Test execution returned non-zero status")
				os.Exit(TestFailed)
			}
		}
	}
}

// Push function pushes the containers to the remote registry
func Push(opts BuildOptions) {
	config := opts.Config

	// If no Git repo exist
	if !isGit() {
		err("No local git repository found, cannot push")
		os.Exit(NoGit)
	}

	if isDirty() {
		err("Git repository has local changes, cannot push")
		os.Exit(GitDirty)
	}

	for _, app := range config.GetApps() {
		for _,branch := range getBranches(opts.All_branches) {
			info("Pushing image %s:%s", app.Image, branch)
			execute("docker", "push", app.Image+":"+branch)
			info("Pushing image %s:%s", app.Image, "latest")
			execute("docker", "push", app.Image+":"+"latest")
		}
	}
}

// Pull function pulls the containers from the remote registry
func Pull(opts BuildOptions) {
	config := opts.Config

	for _, app := range config.GetApps() {
		for _,branch := range getBranches(opts.All_branches) {
			info("Pulling image %s:%s", app.Image, "latest")
			execute("docker", "pull", app.Image+":"+"latest")
			info("Pulling image %s:%s", app.Image, branch)
			execute("docker", "pull", app.Image+":"+branch)
		}
	}
}
