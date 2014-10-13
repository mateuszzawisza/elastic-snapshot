#!/bin/bash
gox -output "./bin/elastic-snapshot_{{.OS}}_{{.Arch}}" -os "linux darwin"
