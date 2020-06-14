#!/bin/sh
export CGO_ENABLED=0
go build -ldflags "-X 'main.BuildTime=`date --iso-8601=seconds`'" "$@"
