#!/usr/bin/env bash

set -eu
set -o pipefail

readonly DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACK_ROOT="$(cd "${DIR}/.." && pwd)"
readonly TEST_CASE="test-case-app"

exit_on_error() {
  exit_code=$1
  last_command=${@:2}
  if [ $exit_code -ne 0 ]; then
    echo >&2 "\"${last_command}\" command failed with exit code ${exit_code}."
    exit $exit_code
  fi
}

echo "Starting Laraboot Integration Test"
echo "WORKING_DIR=$DIR"
echo "BUILDPACK_ROOT=$BUILDPACK_ROOT"
echo "LARABOOT_VERSION=$(laraboot --version)"

laraboot new $TEST_CASE --php-version="8.0.*"
cp -R examples/01-commands/ $TEST_CASE
cd $TEST_CASE

# Copy buildpack to this project as if it was there the whole time
readonly LOCAL_TASK_DIR=".laraboot/tasks/@core/laraboot-commander"
mkdir -p $LOCAL_TASK_DIR

cp "$BUILDPACK_ROOT/buildpack.toml" $LOCAL_TASK_DIR
cp -r "$BUILDPACK_ROOT/bin" $LOCAL_TASK_DIR/bin

# Build the project using this task
laraboot build --pack-params "buildpack $LOCAL_TASK_DIR" -vvv

# Let's check the build commits
# shellcheck disable=SC2046
docker export $(docker ps -lq) -o $TEST_CASE.tar
# extract the container
tar -xf $TEST_CASE.tar
pushd workbench
ls -ltah
git log --pretty=format:"%h - %an, %ar : %s"
echo ""
popd

#See what's on this image
pack inspect $TEST_CASE

#Wait for image to be running
laraboot run --port=9000 -vvv > server.log 2>&1 &
wait-for-it localhost:9000 --timeout=30

# Debug
cat server.log
exit_on_error $? !!