#!/bin/sh
# Fetches the CLDR core.zip used by the table generator (see internal/cldrgen).
#
# Usage: CLDR_VERSION=40 ./fetch-cldr.sh
set -eu

: "${CLDR_VERSION:?CLDR_VERSION must be set, e.g. CLDR_VERSION=40}"

url="https://unicode.org/Public/cldr/${CLDR_VERSION}/core.zip"
out="${PWD}/core.zip"

echo "fetching ${url}"
curl -fL --retry 3 -o "${out}" "${url}"
echo "wrote ${out}"
