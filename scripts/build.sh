#!/bin/bash

version="$1"

ldflags="-s -w -X github.com/ggql/ggql/cli.Version=$version"
target="ggql"

go env -w GOPROXY=https://goproxy.cn,direct

# go tool dist list
CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/$target main.go
