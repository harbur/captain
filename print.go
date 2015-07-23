package captain // import "github.com/harbur/captain"

import "fmt"

func info(text string, arg ...interface{}) {
	text = colorInfo("[") + colorPrefix("CAPTAIN") + colorInfo("]") + " " + text + "\n"
	s := arg
	for i := range s {
		s[i] = colorInfo(s[i])
	}
	fmt.Printf(text, arg...)
}

func err(text string, arg ...interface{}) {
	text = colorErr("[") + colorPrefix("CAPTAIN") + colorErr("]") + " " + text + "\n"
	s := arg
	for i := range s {
		s[i] = colorErr(s[i])
	}
	fmt.Printf(text, s...)
}

func debug(text string, arg ...interface{}) {
	if options.debug {
		text = colorDebug("[") + colorPrefix("CAPTAIN") + colorDebug("]") + " " + text + "\n"
		s := arg
		for i := range s {
			s[i] = colorDebug(s[i])
		}
		fmt.Printf(text, s...)
	}
}
