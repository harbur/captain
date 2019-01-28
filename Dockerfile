# Dockerfile for go-releaser
FROM scratch
COPY mybin /
ENTRYPOINT ["/mybin"]