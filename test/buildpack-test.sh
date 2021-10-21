#!/usr/bin/env bash

set -eu
set -o pipefail

readonly DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACK_ROOT="$(cd "${DIR}/.." && pwd)"
readonly BUILDPACKS_ROOT="$(cd "${BUILDPACK_ROOT}/.." && pwd)"

echo "Init buildpack test"
echo "DIR=$DIR"
echo "BUILDPACK_ROOT=$BUILDPACK_ROOT"
echo "BUILDPACKS_ROOT=$BUILDPACKS_ROOT"

cp -R examples/01-commands/ sample-app
pwd
ls -ltah sample-app

pack build app-name --path sample-app \
  --buildpack paketo-buildpacks/php-dist \
  --buildpack paketo-buildpacks/php-composer \
  --buildpack $BUILDPACK_ROOT \
  --builder paketobuildpacks/builder:full \
  --clear-cache \
  --verbose
