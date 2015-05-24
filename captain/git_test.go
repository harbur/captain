package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitGetRevision(t *testing.T) {
	assert.Equal(t, 7, len(getRevision()), "Git revision has length 7 chars")
}

func TestGitGetBranch(t *testing.T) {
	assert.Equal(t, "master", getBranch(), "Git branch is master")
}

func TestGitIsDirty(t *testing.T) {
	assert.Equal(t, false, isDirty(), "Git does not have local changes")
}

func TestGitIsGit(t *testing.T) {
	assert.Equal(t, true, isGit(), "There is a git repository")
}
