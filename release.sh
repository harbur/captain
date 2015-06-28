#!/usr/bin/env bash

set -e

version=$1

if [ -z "$version" ]; then
  echo "No version passed! Example usage: ./release.sh 1.0.0"
  exit 1
fi

echo "Running tests..."
go test ./...

echo "Update version..."
sed -i.bak 's/fmt\.Println("v[0-9]*\.[0-9]*\.[0-9]*")/fmt.Println("v'$version'")/' captain/cmd.go
rm captain/cmd.go.bak

echo "Build binary..."
make b

echo "Update repository..."
git add captain/cmd.go
git commit -m "Preparing version ${version}"
git tag --message="v$version" "v$version"

echo "v$version tagged."
echo "Now, run 'git push origin master && git push --tags' and publish the release on GitHub."
