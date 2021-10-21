#!/usr/bin/env bash

set -eu
set -o pipefail

readonly DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACK_ROOT="$(cd "${DIR}/.." && pwd)"
readonly BUILDPACKS_ROOT="$(cd "${BUILDPACK_ROOT}/.." && pwd)"

echo "Test"
echo "DIR=$DIR"
echo "BUILDPACK_ROOT=$BUILDPACK_ROOT"
echo "BUILDPACKS_ROOT=$BUILDPACKS_ROOT"

mkdir sample-app
cp -R examples/01-commands/ sample-app

pack build app-name --path sample-app \
  --buildpack paketo-buildpacks/php-dist \
  --buildpack paketo-buildpacks/php-composer \
  --buildpack $BUILDPACK_ROOT \
  --builder paketobuildpacks/builder:full \
  --clear-cache \
  --verbose
