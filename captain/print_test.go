package captain // import "github.com/harbur/captain/captain"

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
	options.debug = true
	debug("test debug %s", "message")
}
