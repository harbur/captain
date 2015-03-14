package captain // import "github.com/harbur/captain"

type StatusError struct {
	error  error
	status int
}

func RealMain() {
	handleCmd()
}
