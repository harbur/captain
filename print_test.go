package captain // import "github.com/harbur/captain"

import (
	"testing"
)

func TestPrintInfo(t *testing.T) {
	info("test info %s", "message")
}

func TestPrintErr(t *testing.T) {
	err("test err %s", "message")
}

func TestPrintDebug(t *testing.T) {
	Debug = true
	defer func() { Debug = false }()
	debug("test debug %s", "message")
}
