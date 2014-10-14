#!/bin/bash
set -e

echo "Running tests..."
./script/test.sh
echo "Done!"
echo

echo "Building..."
gox -output "./bin/elastic-snapshot_{{.OS}}_{{.Arch}}" -os "linux darwin" -arch "amd64 386"
echo "Done!"
