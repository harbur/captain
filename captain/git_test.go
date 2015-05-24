package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitGetRevision(t *testing.T) {
	assert.Equal(t, 7, len(getRevision()), "they should be equal")
}

func TestGitGetBranch(t *testing.T) {
	assert.Equal(t, "master", getBranch(), "they should be equal")
}

func TestGitIsDirty(t *testing.T) {
	assert.Equal(t, false, isDirty(), "they should be equal")
}

func TestGitIsGit(t *testing.T) {
	assert.Equal(t, true, isGit(), "they should be equal")
}
