#!/usr/bin/env bash

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPTSDIR="$(cd "${PROGDIR}/.." && pwd)"
readonly BUILDPACKDIR="$(cd "${SCRIPTSDIR}/.." && pwd)"

echo "PROGDIR=$PROGDIR"
echo "SCRIPTSDIR=$SCRIPTSDIR"
echo "BUILDPACKDIR=$BUILDPACKDIR"

echo "Before build this ..."