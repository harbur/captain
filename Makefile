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
	goconvey
