# Go parameters
GOCMD=go

build:
	docker build -t harbur/captain .

run:
	docker run harbur/captain

b:
	go install -v ./cmd/captain

deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor
	$(GOCMD) mod verify

watch:
	docker run -it --rm --name captain -v "$$PWD":/go/src/github.com/harbur/captain -w /go/src/github.com/harbur/captain golang:1.4 watch -n 1 make b

goconvey:
	goconvey -timeout 10s


cross:
	mkdir -p build
	gox --os windows --os linux --os darwin --arch 386 --arch amd64 github.com/harbur/captain/cmd/captain
	mv captain_{darwin,linux,windows}_* build
