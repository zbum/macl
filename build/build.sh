#!/usr/bin/env bash

env GOOS=darwin GOARCH=amd64 go build -o acl.darwin.amd64 ..
env GOOS=darwin GOARCH=arm64 go build -o acl.darwin.arm64 ..
env GOOS=linux GOARCH=amd64 go build -o acl.linux.amd64 ..
env GOOS=windows GOARCH=amd64 go build -o acl.windows.amd64.exe ..