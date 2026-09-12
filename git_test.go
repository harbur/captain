package captain // import "github.com/harbur/captain"

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitGetRevision(t *testing.T) {
	assert.Equal(t, 7, len(getRevision(false)), "Git revision should have length 7 chars")
}

func TestGitGetRevisionFullSha(t *testing.T) {
	assert.Equal(t, 40, len(getRevision(true)), "Git revision should have a length of 40 chars")
}

// TODO Fails because it assumes current branch is master
func TestGitGetBranch(t *testing.T) {
	// assert.Equal(t, []string{"master"}, getBranches(false), "Git branch should be master")
}

// TODO Fails because it assumes current branch is master
func TestGitGetBranchAllBranches(t *testing.T) {
	// assert.Equal(t, []string{"master"}, getBranches(true), "Git branch should be master")
}

// TODO Fails because vendors/ is not git-ignored.
func TestGitIsDirty(t *testing.T) {
	// assert.Equal(t, false, isDirty(), "Git should not have local changes")
}

func TestGitIsGit(t *testing.T) {
	assert.Equal(t, true, isGit(), "There should be a git repository")
}

// Unlike TestGitIsDirty above, this creates its own controlled untracked file
// instead of relying on the ambient state of the checkout.
func TestGitIsDirtyWithUntrackedFile(t *testing.T) {
	marker := basedir + "/.captain_isdirty_marker_test"
	if err := ioutil.WriteFile(marker, []byte("x"), 0644); err != nil {
		t.Fatalf("failed to create marker file: %v", err)
	}
	defer os.Remove(marker)

	assert.True(t, isDirty(), "An untracked file should make the repository report as dirty")
}
