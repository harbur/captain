package captain // import "github.com/harbur/captain/captain"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageExist(t *testing.T) {
	exist := imageExist("golang","1.4")
    assert.Equal(t, true, exist, "Docker image golang:1.4 should exist")
}

func TestImageDoesNotExist(t *testing.T) {
	exist := imageExist("golang","nonexist")
    assert.Equal(t, false, exist, "Docker image golang:nonexist should not exist")
}
