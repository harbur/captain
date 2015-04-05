package captain // import "github.com/harbur/captain/captain"
import (
	"os"
	"strings"
)

type StatusError struct {
	error  error
	status int
}

func RealMain() {
	handleCmd()
}

func Build(config Config, filter string) {
	var images = config.GetImageNames()

	if filter != "" {
		images = filterImages(images, filter)
	}
	var rev = getRevision()

	for _, value := range images {
		s := strings.Split(value, "=")
		dockerfile, image := s[0], s[1]

		// If no Git repo exist
		if !isGit() {
			// Perfoming [build latest]
			debug("No local git repository found, just building latest")
			// Build latest image
			err := buildImage(dockerfile, image, "latest")
			if err != nil {
				os.Exit(BUILD_FAILED)
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
				err := buildImage(dockerfile, image, "latest")
				if err != nil {
					os.Exit(BUILD_FAILED)
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
						err := tagImage(image, "latest", branch)
						if err != nil {
							os.Exit(TAG_FAILED)
						}
					}
				}
			}
		}
	}
}
