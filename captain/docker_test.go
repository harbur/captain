package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"errors"
	"github.com/stretchr/testify/assert"
)

func TestTagImage(t *testing.T) {
	res := tagImage("golang","1.4","testing")
	assert.Equal(t, nil, res, "Docker tag should not return any error")
}

func TestTagNonexistingImage(t *testing.T) {
	res := tagImage("golang","nonexist","testing")
	expected := errors.New("no such image")
	assert.Equal(t, expected, res, "Docker tag should not return any error")
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
