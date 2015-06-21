package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"os"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestBuildImage(t *testing.T) {
	res := buildImage("test/noCaptainYML/Dockerfile","captain_test","latest")
	assert.Nil(t, res, "Docker build should not return any error")
}

func TestBuildImageError(t *testing.T) {
	res := buildImage("test/noCaptainYML/Dockerfile.error","captain_test","latest")
	assert.NotNil(t,res, "Docker build should return an error")
}

func TestBuildImageCircleCI(t *testing.T) {
	os.Setenv("CIRCLECI", "true")
	res := buildImage("test/noCaptainYML/Dockerfile","captain_test","latest")
	assert.Nil(t, res, "Docker build should not return any error")
}

func TestTagImage(t *testing.T) {
	res := tagImage("golang","1.4","testing")
	assert.Nil(t, res, "Docker tag should not return any error")
}

func TestTagNonexistingImage(t *testing.T) {
	res := tagImage("golang","nonexist","testing")
	expected := errors.New("no such image")
	assert.Equal(t, expected, res, "Docker tag should return an error")
	println()
}

func TestImageExist(t *testing.T) {
	exist := imageExist("golang", "1.4")
	assert.Equal(t, true, exist, "Docker image golang:1.4 should exist")
}

func TestImageDoesNotExist(t *testing.T) {
	exist := imageExist("golang", "nonexist")
	assert.Equal(t, false, exist, "Docker image golang:nonexist should not exist")
}
