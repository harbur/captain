package captain // import "github.com/harbur/captain/captain"

type StatusError struct {
	error  error
	status int
}

func RealMain() {
	handleCmd()
}
