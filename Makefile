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
	cp captain_linux_amd64 captain
	mv captain_darwin_386 captain_Darwin_386
	mv captain_darwin_amd64 captain_Darwin_amd64
	mv captain_darwin_amd64 captain_Darwin_x86_64
	mv captain_linux_386 captain_Linux_386
	mv captain_linux_amd64 captain_Linux_x86_64
	mv captain_windows_386.exe captain_Windows_386.exe
	mv captain_windows_amd64.exe captain_Windows_x86_64.exe
