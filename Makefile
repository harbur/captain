build:
	docker build -t harbur/captain .

run:
	docker run harbur/captain

b:
	go get -v -d github.com/harbur/captain
	go install -v github.com/harbur/captain/cmd/captain

watch:
	docker run -it --rm --name captain -v "$$PWD":/go/src/github.com/harbur/captain -w /go/src/github.com/harbur/captain golang:1.4 watch -n 1 make b

goconvey:
	goconvey -timeout 10s


cross:
	gox --os windows --os linux --os darwin --arch 386 --arch amd64 github.com/harbur/captain/cmd/captain
