package captain // import "github.com/harbur/captain"

import (
	"testing"

	"github.com/harbur/captain/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestGitGetRevision(t *testing.T) {
	assert.Equal(t, 7, len(getRevision()), "Git revision should have length 7 chars")
}

func TestGitGetBranch(t *testing.T) {
	assert.Equal(t, []string{"master"}, getBranches(false), "Git branch should be master")
}

func TestGitGetBranchAllBranches(t *testing.T) {
	assert.Equal(t, []string{"master"}, getBranches(true), "Git branch should be master")
}

func TestGitIsDirty(t *testing.T) {
	assert.Equal(t, false, isDirty(), "Git should not have local changes")
}

func TestGitIsGit(t *testing.T) {
	assert.Equal(t, true, isGit(), "There should be a git repository")
}
