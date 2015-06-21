package captain // import "github.com/harbur/captain/captain"
import (
	"os"
	"sort"
)

// StatusError provides error code and id
type StatusError struct {
	error  error
	status int
}

// RealMain is the Captain entrypoint function
func RealMain() {
	handleCmd()
}

// Pre function executes commands on pre section before build
func Pre(config Config, image string) {
	for _, value := range config.GetPreCommands(image) {
		info("Running pre command: %s", value)
		res := execute("bash", "-c", value)
		if res != nil {
			err("Pre execution returned non-zero status")
			os.Exit(TestFailed)
		}
	}
}

// Post function executes commands on pre section after build
func Post(config Config, image string) {
	for _, value := range config.GetPostCommands(image) {
		info("Running post command: %s", value)
		res := execute("bash", "-c", value)
		if res != nil {
			err("Post execution returned non-zero status")
			os.Exit(TestFailed)
		}
	}
}

// Build function compiles the Containers of the project
func Build(config Config) {
	var images = config.GetImageNames()

	var rev = getRevision()

	// Sort keys to iterate them deterministically
	var keys []string
	for k := range images {
		keys = append(keys, k)
	}
	if (len(keys)==0) {
		err("No Dockerfile(s) found on current or subdirectories, exiting");
		os.Exit(NoDockerfiles)
	}
	sort.Strings(keys)

	for _, dockerfile := range keys {
		Pre(config, dockerfile)
		image := images[dockerfile]
		// If no Git repo exist
		if !isGit() {
			// Perfoming [build latest]
			debug("No local git repository found, just building latest")
			// Build latest image
			res := buildImage(dockerfile, image, "latest")
			if res != nil {
				os.Exit(BuildFailed)
			}

		} else {
			// Skip build if there are no local changes and the commit is already built
			if !isDirty() && imageExist(image, rev) && !options.force {
				// Performing [skip rev|tag rev@latest|tag rev@branch]
				info("Skipping build of %s:%s - image is already built", image, rev)

				// Tag commit image
				tagImage(image, rev, "latest")

				// Tag branch image
				var branch = getBranch()
				switch branch {
				case "HEAD":
					debug("Skipping tag of %s in detached mode", image)
				case "":
					debug("Skipping tag of %s no git repository", image)
				default:
					tagImage(image, rev, branch)
				}

			} else {
				// Performing [build latest|tag latest@rev|tag latest@branch]
				// Build latest image
				res := buildImage(dockerfile, image, "latest")
				if res != nil {
					os.Exit(BuildFailed)
				}
				if isDirty() {
					debug("Skipping tag of %s:%s - local changes exist", image, rev)
				} else {
					// Tag commit image
					tagImage(image, "latest", rev)

					// Tag branch image
					var branch = getBranch()
					switch branch {
					case "HEAD":
						debug("Skipping tag of %s in detached mode", image)
					case "":
						debug("Skipping tag of %s no git repository", image)
					default:
						res := tagImage(image, "latest", branch)
						if res != nil {
							os.Exit(TagFailed)
						}
					}
				}
			}
		}
		Post(config, dockerfile)
	}
}

// Test function executes the tests of the project
func Test(config Config) {
	for _, value := range config.GetUnitTestCommands() {
		info("Running test command: %s", value)
		res := execute("bash", "-c", value)
		if res != nil {
			err("Test execution returned non-zero status")
			os.Exit(TestFailed)
		}
	}
}

// Push function pushes the containers to the remote registry
func Push(config Config) {
	// If no Git repo exist
	if !isGit() {
		err("No local git repository found, cannot push")
		os.Exit(NoGit)
	}

	if isDirty() {
		err("Git repository has local changes, cannot push")
		os.Exit(GitDirty)
	}

	var images = config.GetImageNames()

	// Sort keys to iterate them deterministically
	var keys []string
	for k := range images {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, dockerfile := range keys {
		image := images[dockerfile]
		var branch = getBranch()

		switch branch {
		case "HEAD":
			err("Skipping push of %s in detached mode", image)
		default:
			info("Pushing image %s:%s", image, branch)
			execute("docker", "push", image+":"+branch)
		}
	}
}
