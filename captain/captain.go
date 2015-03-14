package captain

type StatusError struct {
	error  error
	status int
}

func RealMain() {
	handleCmd()
}
