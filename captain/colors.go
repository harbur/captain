package captain // import "github.com/harbur/captain/captain"

import (
	"github.com/fatih/color"
)

var color_prefix = color.New(color.FgWhite, color.Bold).SprintFunc()
var color_info = color.New(color.FgGreen).SprintFunc()
var color_warn = color.New(color.FgYellow).SprintFunc()
var color_err = color.New(color.FgRed).SprintFunc()
