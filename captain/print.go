package captain // import "github.com/harbur/captain/captain"
import "fmt"

func info(text string, arg ...interface{}) {
	text = color_info("[") + color_prefix("CAPTAIN") + color_info("]") + " " + text + "\n"
	s := arg
	for i := range s {
		s[i] = color_info(s[i])
	}
	fmt.Printf(text, arg...)
}

func err(text string, arg ...interface{}) {
	text = color_err("[") + color_prefix("CAPTAIN") + color_err("]") + " " + text + "\n"
	s := arg
	for i := range s {
		s[i] = color_err(s[i])
	}
	fmt.Printf(text, s...)
}
