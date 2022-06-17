#!/bin/sh

# ------------------------------------
# Purpose:
# - Builds executables/binaries (Certificate Client).
#
# Releases:
# - v1.0.0 - 2022-06-17: initial release
#
# Remarks:
# - go tool dist list
# ------------------------------------

# set -o xtrace
set -o verbose

# compile 'aix'
env GOOS=aix GOARCH=ppc64 go build -o build/aix-ppc64/osmnotes2csv

# compile 'darwin'
env GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/osmnotes2csv
env GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/osmnotes2csv

# compile 'dragonfly'
env GOOS=dragonfly GOARCH=amd64 go build -o build/dragonfly-amd64/osmnotes2csv

# compile 'freebsd'
env GOOS=freebsd GOARCH=amd64 go build -o build/freebsd-amd64/osmnotes2csv
env GOOS=freebsd GOARCH=arm64 go build -o build/freebsd-arm64/osmnotes2csv

# compile 'illumos'
env GOOS=illumos GOARCH=amd64 go build -o build/illumos-amd64/osmnotes2csv

# compile 'linux'
env GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/osmnotes2csv
env GOOS=linux GOARCH=arm64 go build -o build/linux-arm64/osmnotes2csv
env GOOS=linux GOARCH=mips64 go build -o build/linux-mips64/osmnotes2csv
env GOOS=linux GOARCH=mips64le go build -o build/linux-mips64le/osmnotes2csv
env GOOS=linux GOARCH=ppc64 go build -o build/linux-ppc64/osmnotes2csv
env GOOS=linux GOARCH=ppc64le go build -o build/linux-ppc64le/osmnotes2csv
env GOOS=linux GOARCH=riscv64 go build -o build/linux-riscv64/osmnotes2csv
env GOOS=linux GOARCH=s390x go build -o build/linux-s390x/osmnotes2csv

# compile 'netbsd'
env GOOS=netbsd GOARCH=amd64 go build -o build/netbsd-amd64/osmnotes2csv
env GOOS=netbsd GOARCH=arm64 go build -o build/netbsd-arm64/osmnotes2csv

# compile 'openbsd'
env GOOS=openbsd GOARCH=amd64 go build -o build/openbsd-amd64/osmnotes2csv
env GOOS=openbsd GOARCH=arm64 go build -o build/openbsd-arm64/osmnotes2csv
env GOOS=openbsd GOARCH=mips64 go build -o build/openbsd-mips64/osmnotes2csv

# compile 'solaris'
env GOOS=solaris GOARCH=amd64 go build -o build/solaris-amd64/osmnotes2csv

# compile 'windows'
env GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/osmnotes2csv.exe
env GOOS=windows GOARCH=386 go build -o build/windows-386/osmnotes2csv.exe
env GOOS=windows GOARCH=arm go build -o build/windows-arm/osmnotes2csv.exe
