#!/bin/sh

mkdir -p .releases

GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o .releases/hybr hybr/cmd/cli && upx .releases/hybr
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o .releases/hybr-server hybr/cmd/server && upx .releases/hybr-server
