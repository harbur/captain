package captain // import "github.com/harbur/captain"

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Fixtures
var validApp = App{
	Image: "",
	Pre:   []string{"echo running pre"},
	Post:  []string{"echo running post"},
}

var invalidApp = App{
	Pre:  []string{"nonexistingPreCommand"},
	Post: []string{"nonexistingPostCommand"},
}

// Pre Command
func TestPre(t *testing.T) {
	res := Pre(validApp)
	assert.Nil(t, res, "No error returned")
}

func TestPreFail(t *testing.T) {
	res := Pre(invalidApp)
	assert.NotNil(t, res, "Error returned")
}

// Post Command
func TestPost(t *testing.T) {
	res := Post(validApp)
	assert.Nil(t, res, "No error returned")
}

func TestPostFail(t *testing.T) {
	res := Post(invalidApp)
	assert.NotNil(t, res, "Error returned")
}

// Build Command
func TestBuild(t *testing.T) {
	var testConfig = readConfig(configFile(basedir + "/test/Simple/captain.yml"))

	var buildOpts = BuildOptions{
		Config: testConfig,
	}

	Build(buildOpts)
}

// Test Command
func TestTest(t *testing.T) {
	var testConfig = readConfig(configFile(basedir + "/test/Simple/captain.yml"))

	var buildOpts = BuildOptions{
		Config: testConfig,
	}

	Test(buildOpts)
}

// Pull Command
func TestPullNoBranchTags(t *testing.T) {
	var testConfig = readConfig(configFile(basedir + "/test/alpine/captain.yml"))

	var buildOpts = BuildOptions{
		Config:      testConfig,
		Branch_tags: false,
	}
	Pull(buildOpts)
}

// Purge Command
func TestPurge(t *testing.T) {
	var testConfig = readConfig(configFile(basedir + "/test/alpine/captain.yml"))

	var buildOpts = BuildOptions{
		Config: testConfig,
	}
	Purge(buildOpts)
}

// SelfUpdate Command
func TestSelfUpdate(t *testing.T) {
	// First Time Self update
	SelfUpdate()
	// Already Installed last version
	SelfUpdate()
}

func TestDownloadFile(t *testing.T) {
	res := downloadFile("/tmp/captain.html", "https://github.com/harbur/captain")
	assert.Nil(t, res, "captain")
}

func TestFindLastVersion(t *testing.T) {
	res := findLastVersion()
	assert.NotNil(t, res, "Last version exists")
}

// captureStdout temporarily redirects os.Stdout while f runs, returning
// everything written to it. Used to observe Build's info() logging, which is
// the only externally-visible signal of whether a build was skipped.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w

	f()

	os.Stdout = old
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// scratchApp/scratchConfig build a minimal, network-free App/Config backed by
// test/scratch/Dockerfile (FROM scratch + a COPY), so Build's real docker
// build/tag logic can be exercised without needing to reach a registry.
func scratchApp(image string) App {
	return App{Build: basedir + "/test/scratch/Dockerfile", Image: image}
}

func scratchConfig(image string) Config {
	yml := "scratchapp:\n  build: " + basedir + "/test/scratch/Dockerfile\n  image: " + image + "\n"
	return unmarshal([]byte(yml))
}

// Build Command - fresh image, never built before at this revision
func TestBuildFreshImage(t *testing.T) {
	image := "captain_test_build_fresh"
	buildOpts := BuildOptions{Config: scratchConfig(image)}

	output := captureStdout(t, func() {
		Build(buildOpts)
	})

	assert.Contains(t, output, "Building image", "a never-built image should trigger an actual build")

	app := scratchApp(image)
	assert.True(t, imageExist(app, "latest"), "latest tag should exist after build")
	assert.True(t, imageExist(app, getRevision(false)), "revision tag should exist after building a clean repo")
}

// Build Command - user-supplied --tag should be applied on top of a fresh build
func TestBuildWithUserTag(t *testing.T) {
	image := "captain_test_build_usertag"
	buildOpts := BuildOptions{Config: scratchConfig(image), Tag: "mytag"}

	captureStdout(t, func() {
		Build(buildOpts)
	})

	app := scratchApp(image)
	assert.True(t, imageExist(app, "mytag"), "user-provided tag should be applied after build")
}

// Build Command - a revision already built (e.g. re-run of CI on the same commit)
// should be re-tagged rather than rebuilt.
func TestBuildSkipsWhenAlreadyBuiltForRevision(t *testing.T) {
	image := "captain_test_build_skip"
	app := scratchApp(image)
	rev := getRevision(false)

	// Simulate that this revision was already built by a previous run.
	if res := buildImage(app, "latest", false); res != nil {
		t.Fatalf("setup build failed: %v", res)
	}
	if res := tagImage(app, "latest", rev); res != nil {
		t.Fatalf("setup tag failed: %v", res)
	}

	buildOpts := BuildOptions{Config: scratchConfig(image), Tag: "skiptag"}
	output := captureStdout(t, func() {
		Build(buildOpts)
	})

	assert.NotContains(t, output, "Building image", "build should be skipped when the revision is already built")
	assert.True(t, imageExist(app, "latest"), "latest should still be (re-)tagged from the existing revision image")
	assert.True(t, imageExist(app, "skiptag"), "user tag should still apply even when the build itself is skipped")
}

// Build Command - --force should rebuild even when the revision was already built.
func TestBuildForceRebuildsEvenWhenAlreadyBuilt(t *testing.T) {
	image := "captain_test_build_force"
	app := scratchApp(image)
	rev := getRevision(false)

	if res := buildImage(app, "latest", false); res != nil {
		t.Fatalf("setup build failed: %v", res)
	}
	if res := tagImage(app, "latest", rev); res != nil {
		t.Fatalf("setup tag failed: %v", res)
	}

	buildOpts := BuildOptions{Config: scratchConfig(image), Force: true}
	output := captureStdout(t, func() {
		Build(buildOpts)
	})

	assert.Contains(t, output, "Building image", "force should trigger a rebuild even when the revision is already built")
}

// Build Command - local changes should suppress revision/branch tagging, since
// the built image cannot be trusted to represent any particular commit.
func TestBuildSkipsTaggingWhenRepoIsDirty(t *testing.T) {
	image := "captain_test_build_dirty"
	marker := basedir + "/.captain_dirty_marker_test"
	if err := ioutil.WriteFile(marker, []byte("dirty"), 0644); err != nil {
		t.Fatalf("failed to create dirty marker: %v", err)
	}
	defer os.Remove(marker)

	assert.True(t, isDirty(), "repository should be considered dirty with an untracked file present")

	buildOpts := BuildOptions{Config: scratchConfig(image)}
	captureStdout(t, func() {
		Build(buildOpts)
	})

	app := scratchApp(image)
	assert.True(t, imageExist(app, "latest"), "latest should still be tagged directly by the build itself")
	assert.False(t, imageExist(app, getRevision(false)), "revision tag should be skipped while the repository has local changes")
}
