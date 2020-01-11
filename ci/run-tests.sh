#!/usr/bin/env sh

# Run virtual frame buffer with Display :99
Xvfb :99 -ac -screen 0 "1920x1080x24" -nolisten tcp &

go test -race -v ./...
go test -run "TestWindows_Open" -v ./opengl
go test -run "TestWindow_Zoom" -v ./opengl
