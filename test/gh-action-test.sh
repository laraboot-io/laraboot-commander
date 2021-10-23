#!/usr/bin/env bash

set -eu
set -o pipefail

readonly DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Init gh-action test"
echo "DIR=$DIR"

cp -R examples/01-commands/ gh-app

# Here we're testing our previously created buildpack
pack build app-name --path gh-app \
  --buildpack docker://my-buildpack \
  --builder paketobuildpacks/builder:full \
  --clear-cache \
  --verbose
