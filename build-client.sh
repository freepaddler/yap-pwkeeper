#!/bin/sh

usage() {
    cat <<EOF
    This script builds pwkeeper client for different platforms.
    Script expects exactly one argument - platform to build for.
    Supported platforms:
        darwin_amd64, darwin_arm64
        linux_amd64, linux_arm64
        windows_amd64
    Special platform name "current" build for current architecture

EOF
    exit 1
}

[ $# -ne 1 ] && usage

case $1 in
    darwin_amd64)
        goos="darwin"
        goarch="amd64"
        client="client_${goos}_${goarch}"
        ;;
    darwin_arm64)
        goos="darwin"
        goarch="arm64"
        client="client_${goos}_${goarch}"
        ;;
    linux_amd64)
        goos="linux"
        goarch="amd64"
        client="client_${goos}_${goarch}"
        ;;
    linux_arm64)
        goos="linux"
        goarch="arm64"
        client="client_${goos}_${goarch}"
        ;;
    windows_amd64)
        goos="windows"
        goarch="amd64"
        client="client_${goos}_${goarch}.exe"
        ;;
    current)
        client="client"
        ;;
    *)
        usage
        ;;
esac

env GOOS=$goos GOARCH=$goarch  go build -ldflags \
"-X 'main.buildVersion=$(git describe --tag --always 2>/dev/null)' \
-X 'main.buildDate=$(date)'" \
-o $client cmd/client/client.go