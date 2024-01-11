#!/bin/sh

go build -ldflags \
"-X 'main.buildVersion=$(git describe --tag --always 2>/dev/null)' \
-X 'main.buildDate=$(date)'" \
-o server cmd/server/server.go