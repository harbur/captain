package captain // import "github.com/harbur/captain"

import (
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// Debug can be turned on to enable debug mode.
var Debug bool

// StatusError provides error code and id
type StatusError struct {
	error  error
	status int
}

// Pre function executes commands on pre section before build
func Pre(app App) error {
	for _, value := range app.Pre {
		info("Running pre command: %s", value)
		res := execute("bash", "-c", value)
		if res != nil {
			return res
		}
	}
	return nil
}

// Post function executes commands on pre section after build
func Post(app App) error {
	for _, value := range app.Post {
		info("Running post command: %s", value)
		res := execute("bash", "-c", value)
		if res != nil {
			return res
		}
	}
	return nil
}

type BuildOptions struct {
	Config       Config
	Tag          string
	Force        bool
	All_branches bool
	Long_sha     bool
	Branch_tags  bool
	Commit_tags  bool
}

// Build function compiles the Containers of the project
func Build(opts BuildOptions) {
	config := opts.Config

	var rev = getRevision(opts.Long_sha)

	// For each App
	for _, app := range config.GetApps() {
		// If no Git repo exist
		if !isGit() {
			// Perfoming [build latest]
			debug("No local git repository found, just building latest")

			// Execute Pre commands
			if res := Pre(app); res != nil {
				err("Pre execution returned non-zero status")
				return
			}

			// Build latest image
			res := buildImage(app, "latest", opts.Force)
			if res != nil {
				os.Exit(BuildFailed)
			}

			// Add additional user-defined Tag
			if opts.Tag != "" {
				tagImage(app, "latest", opts.Tag)
			}
		} else {
			// Skip build if there are no local changes and the commit is already built
			if !isDirty() && imageExist(app, rev) && !opts.Force {
				// Performing [skip rev|tag rev@latest|tag rev@branch]
				info("Skipping build of %s:%s - image is already built", app.Image, rev)

				// Tag commit image
				tagImage(app, rev, "latest")

				// Tag branch image
				for _, branch := range getBranches(opts.All_branches) {
					res := tagImage(app, rev, branch)
					if res != nil {
						os.Exit(TagFailed)
					}
				}

				// Add additional user-defined Tag
				if opts.Tag != "" {
					tagImage(app, rev, opts.Tag)
				}
			} else {
				// Performing [build latest|tag latest@rev|tag latest@branch]

				// Execute Pre commands
				if res := Pre(app); res != nil {
					err("Pre execution returned non-zero status")
				}

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
					for _, branch := range getBranches(opts.All_branches) {
						res := tagImage(app, "latest", branch)
						if res != nil {
							os.Exit(TagFailed)
						}
					}

					// Add additional user-defined Tag
					if opts.Tag != "" {
						tagImage(app, rev, opts.Tag)
					}
				}
			}
		}

		// Execute Post commands
		if res := Post(app); res != nil {
			err("Post execution returned non-zero status")
		}
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
				return
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
		for _, branch := range getBranches(opts.All_branches) {
			info("Pushing image %s:%s", app.Image, "latest")
			if res := pushImage(app.Image, "latest"); res != nil {
				err("Push returned non-zero status")
				os.Exit(ExecuteFailed)
			}
			if opts.Branch_tags {
				info("Pushing image %s:%s", app.Image, branch)
				if res := pushImage(app.Image, branch); res != nil {
					err("Push returned non-zero status")
					os.Exit(ExecuteFailed)
				}
			}
			if opts.Commit_tags {
				rev := getRevision(opts.Long_sha)
				info("Pushing image %s:%s", app.Image, rev)
				if res := pushImage(app.Image, rev); res != nil {
					err("Push returned non-zero status")
					os.Exit(ExecuteFailed)
				}
			}

			// Add additional user-defined Tag
			if opts.Tag != "" {
				info("Pushing image %s:%s", app.Image, opts.Tag)
				if res := pushImage(app.Image, opts.Tag); res != nil {
					err("Push returned non-zero status")
					os.Exit(ExecuteFailed)
				}
			}
		}
	}
}

// Pull function pulls the containers from the remote registry
func Pull(opts BuildOptions) {
	config := opts.Config

	for _, app := range config.GetApps() {
		for _, branch := range getBranches(opts.All_branches) {
			info("Pulling image %s:%s", app.Image, "latest")
			if res := pullImage(app.Image, "latest"); res != nil {
				err("Pull returned non-zero status")
				os.Exit(ExecuteFailed)
			}
			if opts.Branch_tags {
				info("Pulling image %s:%s", app.Image, branch)
				if res := pullImage(app.Image, branch); res != nil {
					err("Pull returned non-zero status")
					os.Exit(ExecuteFailed)
				}
			}
			if opts.Commit_tags {
				rev := getRevision(opts.Long_sha)
				info("Pulling image %s:%s", app.Image, rev)
				if res := pullImage(app.Image, rev); res != nil {
					err("Pull returned non-zero status")
					os.Exit(ExecuteFailed)
				}
			}

			// Add additional user-defined Tag
			if opts.Tag != "" {
				info("Pulling image %s:%s", app.Image, opts.Tag)
				if res := pullImage(app.Image, opts.Tag); res != nil {
					err("Pull returned non-zero status")
					os.Exit(ExecuteFailed)
				}
			}
		}
	}
}

// Purge function purges the stale images
func Purge(opts BuildOptions) {
	config := opts.Config

	// For each App
	for _, app := range config.GetApps() {
		var tags = []string{}

		// Retrieve the list of the existing Image tags
		for _, img := range getImages(app) {
			tags = append(tags, img.RepoTags...)
		}

		// Remove from the list: The latest image
		for key, tag := range tags {
			if tag == app.Image+":latest" {
				tags = append(tags[:key], tags[key+1:]...)
			}
		}

		// Remove from the list: The current commit-id
		for key, tag := range tags {
			if tag == app.Image+":"+getRevision(opts.Long_sha) {
				tags = append(tags[:key], tags[key+1:]...)
			}
		}

		// Remove from the list: The working-dir git branches
		for _, branch := range getBranches(opts.All_branches) {
			for key, tag := range tags {
				if tag == app.Image+":"+branch {
					tags = append(tags[:key], tags[key+1:]...)
				}
			}
		}

		// Proceed with deletion of Images
		for _, tag := range tags {
			info("Deleting image %s", tag)
			res := removeImage(tag)
			if res != nil {
				err("Deleting image failed: %s", res)
				os.Exit(DeleteImageFailed)
			}
		}
	}
}

func SelfUpdate() {
	captainDir := filepath.FromSlash(os.Getenv("HOME") + "/.captain")
	binariesDir := filepath.FromSlash(captainDir + "/binaries")
	binDir := filepath.FromSlash(captainDir + "/bin")

	kernel := runtime.GOOS
	arch := runtime.GOARCH
	captainSymlinkPath := filepath.FromSlash(binDir + "/captain")
	currentVersionPath, _ := os.Readlink(captainSymlinkPath)

	info("Checking the last version of Captain...")
	version := findLastVersion()
	downloadUrl := fmt.Sprintf("https://github.com/harbur/captain/releases/download/%s/captain_%s_%s", version, kernel, arch)
	downloadedVersionPath := filepath.FromSlash(binariesDir + "/captain-" + version)

	if currentVersionPath == downloadedVersionPath {
		info("You have installed the last version of captain (%s)", version)
		return
	}

	info("New version available, start downloading %s", version)

	// create binaries dir
	if err := os.MkdirAll(binariesDir, 0755); err != nil {
		fmt.Println(err)
	}

	// download new version
	if err := downloadFile(downloadedVersionPath, downloadUrl); err != nil {
		fmt.Println(err)
	}

	// grant excution to download version
	if err := os.Chmod(downloadedVersionPath, 0755); err != nil {
		fmt.Println(err)
	}

	info("Downloaded captain %s", version)

	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Println(err)
	}

	if _, err := os.Stat(captainSymlinkPath); err == nil {
		os.Remove(captainSymlinkPath)
	}

	if err := os.Symlink(downloadedVersionPath, captainSymlinkPath); err != nil {
		fmt.Println(err)
	}

	info("Installed captain %s", version)
}

func findLastVersion() string {
	url := "https://raw.githubusercontent.com/harbur/captain/master/VERSION"

	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func downloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// create and start bar
	bar := pb.New64(resp.ContentLength).SetUnits(pb.U_BYTES).Start()
	defer bar.Finish()

	// create proxy reader
	reader := bar.NewProxyReader(resp.Body)

	// and copy from pb reader
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	return nil
}
