#!/usr/bin/env bash

set -eu
set -o pipefail

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACKDIR="$(cd "${PROGDIR}/.." && pwd)"

function main() {
  while [[ "${#}" != 0 ]]; do
    case "${1}" in
    --help | -h)
      shift 1
      usage
      exit 0
      ;;

    "")
      # skip if the argument is empty
      shift 1
      ;;

    *)
      util::print::error "unknown argument \"${1}\""
      ;;
    esac
  done

  cmd::build
}

function usage() {
  cat <<-USAGE
release-buildpacks-aws.sh [OPTIONS]

Release the buildpack into AWS s3 bucket.

OPTIONS
  --help  -h  prints the command usage
USAGE
}

function cmd::build() {

  : ${STAGE:=dev}

  aws s3 sync dist s3://buildpacks.laraboot.io/core/${STAGE}/
}

main "${@:-}"
