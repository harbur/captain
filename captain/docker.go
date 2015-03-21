package captain // import "github.com/harbur/captain/captain"

func buildImage(dockerfile string, image string) {
	info("Building image %s", image)
	execute("docker", "build", "-f", dockerfile, "-t", image, ".")
}

func tagImage(origin string, target string) {
	info("Tagging image as %s", target)
	execute("docker", "tag", "-f", origin, target)
}
