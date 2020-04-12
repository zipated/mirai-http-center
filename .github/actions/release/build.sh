#!/bin/bash

go mod download

GO_VERSION=$(go version)
BUILD_TIME=$(date)
ldflags="-X 'main.version=${VERSION}' -X 'main.goVersion=${GO_VERSION}' -X 'main.buildTime=${BUILD_TIME}' -X 'main.commitID=${COMMIT_ID}'"

CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o mirai-http-center_linux_32 -ldflags "${ldflags}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mirai-http-center_linux_64 -ldflags "${ldflags}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o mirai-http-center_linux_armv6 -ldflags "${ldflags}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o mirai-http-center_linux_armv7 -ldflags "${ldflags}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o mirai-http-center_linux_arm64 -ldflags "${ldflags}"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o mirai-http-center_macos -ldflags "${ldflags}"
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o mirai-http-center_windows_32.exe -ldflags "${ldflags}"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o mirai-http-center_windows_64.exe -ldflags "${ldflags}"
