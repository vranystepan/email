#!/bin/sh

# make sure SERVICE is not propagated from the parent
unset SERVICE

# set Go configuration
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

# export Go configuration
export CGO_ENABLED GOOS GOARCH

while getopts 's:' c
do
  case $c in
    s) SERVICE=$OPTARG ;;
  esac
done

[ -z "${SERVICE}" ] && { echo "SERVICE (-s) can't be empty!"; exit 1; }
echo "building ${SERVICE}"

mkdir -p ./bin/${SERVICE}
go build -o ./bin/${SERVICE}/main ./cmd/${SERVICE}/main.go
ls -lah ./bin/${SERVICE}/main
