#!/bin/sh
go build -ldflags "-X 'main.BuildTime=`date --iso-8601=seconds`'" "$@"
