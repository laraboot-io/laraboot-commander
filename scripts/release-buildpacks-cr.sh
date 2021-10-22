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
release-buildpacks-cr.sh [OPTIONS]

Release the buildpack into a Container Registry

OPTIONS
  --help  -h  prints the command usage
USAGE
}

function cmd::build() {

  : ${IMAGE_TAG:=dev}

  readonly name="laraboot-rector"
  readonly buildpack_id="laraboot-buildpacks/laraboot-rector"

  echo "  ----> Id: $buildpack_id"
  echo "  ----> Tag: $ECR_REGISTRY/$buildpack_id:$IMAGE_TAG"
  docker tag $name $ECR_REGISTRY/$buildpack_id:$IMAGE_TAG
  docker push $ECR_REGISTRY/$buildpack_id:"$IMAGE_TAG"
}

main "${@:-}"
