package captain // import "github.com/harbur/captain"

import (
	"io/ioutil"
	"log"
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
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defer os.Chdir(pwd)
	dir, err := ioutil.TempDir("", "captain_tmp-")
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(dir)
	defer os.RemoveAll(dir)
	CreateRepoWithCommits()
	// when tag checked out
	oneliner("git", "checkout", "v1.0")
	assert.Equal(
		t,
		[]string{"b-three", "v1.0", "v1.1"},
		getBranches(false),
		"Git branch should be those pointing to v1.0")
	// when branch checked out
	oneliner("git", "checkout", "b-two")
	assert.Equal(
		t,
		[]string{"b-three", "v1.0", "v1.1"}, // returns b-three instead of b-two since it is first alphanumerically
		getBranches(false),
		"Git branch should be those pointing to b-two")
	// when branch with no tags checked out
	oneliner("git", "checkout", "b-one")
	assert.Equal(
		t,
		[]string{"b-one"},
		getBranches(false),
		"Git branch should be those pointing to b-one")
	// when tag with no branches checked out
	oneliner("git", "checkout", "v3.0")
	assert.Equal(
		t,
		[]string{"undefined", "v3.0"}, // "undefined" is a string returned by git. Tempted to leave as is.
		getBranches(false),
		"Git branch should be those pointing to v3.0")
}

func TestGitGetBranchAllBranches(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defer os.Chdir(pwd)
	dir, err := ioutil.TempDir("", "captain_tmp-")
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(dir)
	defer os.RemoveAll(dir)
    CreateRepoWithCommits()
	// when tag checked out
	oneliner("git", "checkout", "v1.0")
	assert.Equal(
		t,
		[]string{"b-three", "b-two", "v1.0", "v1.1"},
	   getBranches(true),
	   "Git branch should be those pointing to v1.0")
	// when branch checked out
	oneliner("git", "checkout", "b-two")
	assert.Equal(
		t,
		[]string{"b-three", "b-two", "v1.0", "v1.1"},
	   getBranches(true),
	   "Git branch should be those pointing to b-two")
	// when branch with no tags checked out
	oneliner("git", "checkout", "b-one")
	assert.Equal(
		t,
		[]string{"b-one", "master"},
	   getBranches(true),
	   "Git branch should be those pointing to b-one")
	// when tag with no branches checked out
	oneliner("git", "checkout", "v3.0")
	assert.Equal(
		t,
		[]string{"v3.0"},
	    getBranches(true),
	    "Git branch should be those pointing to v3.0")
}

// TODO Fails because vendors/ is not git-ignored.
func TestGitIsDirty(t *testing.T) {
	// assert.Equal(t, false, isDirty(), "Git should not have local changes")
}

func TestGitIsGit(t *testing.T) {
	assert.Equal(t, true, isGit(), "There should be a git repository")
}

func CreateRepoWithCommits(){
	// create a git repository of the following form
	//   commit-A   => master, b-one
	//   commit-B   => b-two, b-three, v1.0, v1.1
	//   commit-C   => b-four, v2.0
	//   commit-D   => v3.0
	oneliner("git", "init")
	CommitEmptyFile("a.txt")
	oneliner("git", "checkout", "-b", "b-one")
	oneliner("git", "checkout", "-b", "b-two")
	CommitEmptyFile("b.txt")
	oneliner("git", "checkout", "-b", "b-three")
	oneliner("git", "tag", "-a", "v1.0", "-m", "1.0")
	oneliner("git", "tag", "v1.1")
	oneliner("git", "checkout", "-b", "b-four")
	CommitEmptyFile("c.txt")
	oneliner("git", "tag", "-a", "v2.0", "-m", "2.0")
	oneliner("git", "checkout", "-b", "b-five")
	CommitEmptyFile("d.txt")
	oneliner("git", "tag", "-a", "v3.0", "-m", "3.0")
	oneliner("git", "checkout", "v3.0")
	oneliner("git", "branch", "-D", "b-five")
}

func CommitEmptyFile(name string){
	CreateEmptyFile(name)
	oneliner("git", "add", ".")
	oneliner("git", "commit", "-m", "wip")
}

func CreateEmptyFile(name string){
	emptyFile, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer emptyFile.Close()
}
