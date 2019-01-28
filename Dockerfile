FROM golang:1.11 as build
WORKDIR /go/src/github.com/harbur/captain
COPY . .
RUN curl https://glide.sh/get | sh
RUN glide install 
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /captain ./cmd/captain/
RUN chmod +x /captain

FROM alpine:3.8

COPY --from=build /captain /bin/captain
ENTRYPOINT [ "/bin/captain" ]
