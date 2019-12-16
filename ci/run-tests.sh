#!/usr/bin/env sh

# Run virtual frame buffer with Display :99
Xvfb :99 -ac -screen 0 "1920x1080x24" -nolisten tcp &

# TODO go test -race does not work on Alpine :(
go test -v ./...
