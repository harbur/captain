package captain // import "github.com/harbur/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitGetRevision(t *testing.T) {
	assert.Equal(t, 7, len(getRevision(false)), "Git revision should have length 7 chars")
}

func TestGitGetRevisionFullSha(t *testing.T) {
	assert.Equal(t, 40, len(getRevision(true)), "Git revision should have a length of 40 chars")
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
