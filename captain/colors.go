package captain

import (
	"github.com/fatih/color"
)

var info = color.New(color.FgBlue).SprintFunc()
var warn = color.New(color.FgYellow).SprintFunc()
var err = color.New(color.FgRed).SprintFunc()
