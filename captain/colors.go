package captain

import (
	"github.com/fatih/color"
)

var prefix = color.New(color.FgWhite, color.Bold).SprintFunc()
var info = color.New(color.FgBlue).SprintFunc()
var warn = color.New(color.FgYellow).SprintFunc()
var err = color.New(color.FgRed).SprintFunc()
