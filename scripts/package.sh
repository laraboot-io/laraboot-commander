#!/usr/bin/env bash

mkdir dist

go build -ldflags="-s -w" -o ./bin/detect ./cmd/detect/main.go && \
go build -ldflags="-s -w" -o ./bin/build ./cmd/build/main.go

# https://buildpacks.io/docs/buildpack-author-guide/package-a-buildpack/
pack buildpack package dist/laraboot-rector.cnb --config ./package.toml --format file
pack buildpack package laraboot-commander --config ./package.toml