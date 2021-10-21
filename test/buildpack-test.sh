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
cp examples/01-commands/laraboot.json sample-app/laraboot.json
echo "<?php echo 1; " > sample-app/server.php

pack build app-name --path $BUILDPACK_ROOT/sample-app \
  --buildpack paketo-buildpacks/php-dist \
  --buildpack paketo-buildpacks/php-composer \
  --buildpack $BUILDPACK_ROOT \
  --builder paketobuildpacks/builder:full \
  --clear-cache \
  --verbose
