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

curl -s "https://laravel.build/sample-app" | bash
cp assets/rector.php sample-app/rector.php
cp assets/laraboot.json sample-app/laraboot.json
cp assets/buildpack.yml sample-app/buildpack.yml

# until https://github.com/paketo-buildpacks/php-dist/issues/201 is fixed
pushd sample-app
cat composer.json | jq -r ".require.php=\"8.0.*\"" >composer.tmp &&
  mv composer.tmp composer.json &&
  rm composer.lock
popd

pack build app-name --path $BUILDPACK_ROOT/sample-app \
  --buildpack paketo-buildpacks/php-dist \
  --buildpack paketo-buildpacks/php-composer \
  --buildpack $BUILDPACK_ROOT \
  --builder paketobuildpacks/builder:full \
  --clear-cache \
  --verbose
