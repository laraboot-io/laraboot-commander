#!/usr/bin/env bash
set -eu
set -o pipefail

WORKSPACE_DIR="${WORKSPACE_DIR:/github/workspace}"

pwd

go build -ldflags="-s -w" -o ./bin/detect ./cmd/detect/main.go &&
go build -ldflags="-s -w" -o ./bin/build ./cmd/build/main.go &&
pack buildpack package "$WORKSPACE_DIR/dist/laraboot-commander.cnb" --config ./package.toml --format file

exit 0
