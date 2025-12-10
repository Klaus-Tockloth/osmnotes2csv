#!/bin/sh

# ------------------------------------
# Purpose:
# - Builds executables/binaries (Certificate Client).
#
# Releases:
# - v1.0.0 - 2022-06-17: initial release
# - v1.1.0 - 2025-12-09: revised
#
# Remarks:
# - go tool dist list
# ------------------------------------

# set -o xtrace
set -o verbose
# set -o errexit

# linting
golangci-lint run --no-config --enable gocritic
revive

# security checks
govulncheck ./...
gosec -exclude=G304 ./...

# compile 'darwin'
env GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/osmnotes2csv
env GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/osmnotes2csv

# compile 'linux'
env GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/osmnotes2csv
env GOOS=linux GOARCH=arm64 go build -o build/linux-arm64/osmnotes2csv

# compile 'windows'
env GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/osmnotes2csv.exe
env GOOS=windows GOARCH=arm go build -o build/windows-arm/osmnotes2csv.exe
