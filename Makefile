build:
	docker build -t harbur/captain .

run:
	docker run harbur/captain

b:
	go get -v -d github.com/harbur/captain
	go install -v github.com/harbur/captain

watch:
	docker run -it --rm --name captain -v "$$PWD":/go/src/github.com/harbur/captain -w /go/src/github.com/harbur/captain golang:1.4 watch -n 1 make b

goconvey:
	goconvey -timeout 10s


cross:
	docker run --rm -v "$$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=windows -e GOARCH=386 golang:1.4.2-cross sh -c 'go get ./...; go build -v'
