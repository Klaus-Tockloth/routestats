#!/bin/sh

# ------------------------------------
# Purpose:
# - Builds executables/binaries.
#
# Releases:
# - v1.0.0 - 2023-10-10: initial release
#
# Remarks:
# - go tool dist list
# ------------------------------------

# set -o xtrace
set -o verbose

# compile 'aix'
env GOOS=aix GOARCH=ppc64 go build -o build/aix-ppc64/routestats

# compile 'darwin'
env GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/routestats
env GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/routestats

# compile 'dragonfly'
env GOOS=dragonfly GOARCH=amd64 go build -o build/dragonfly-amd64/routestats

# compile 'freebsd'
env GOOS=freebsd GOARCH=amd64 go build -o build/freebsd-amd64/routestats
env GOOS=freebsd GOARCH=arm64 go build -o build/freebsd-arm64/routestats

# compile 'illumos'
env GOOS=illumos GOARCH=amd64 go build -o build/illumos-amd64/routestats

# compile 'linux'
env GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/routestats
env GOOS=linux GOARCH=arm64 go build -o build/linux-arm64/routestats
env GOOS=linux GOARCH=mips64 go build -o build/linux-mips64/routestats
env GOOS=linux GOARCH=mips64le go build -o build/linux-mips64le/routestats
env GOOS=linux GOARCH=ppc64 go build -o build/linux-ppc64/routestats
env GOOS=linux GOARCH=ppc64le go build -o build/linux-ppc64le/routestats
env GOOS=linux GOARCH=riscv64 go build -o build/linux-riscv64/routestats

# compile 'netbsd'
env GOOS=netbsd GOARCH=amd64 go build -o build/netbsd-amd64/routestats
env GOOS=netbsd GOARCH=arm64 go build -o build/netbsd-arm64/routestats

# compile 'openbsd'
env GOOS=openbsd GOARCH=amd64 go build -o build/openbsd-amd64/routestats
env GOOS=openbsd GOARCH=arm64 go build -o build/openbsd-arm64/routestats
env GOOS=openbsd GOARCH=mips64 go build -o build/openbsd-mips64/routestats

# compile 'solaris'
env GOOS=solaris GOARCH=amd64 go build -o build/solaris-amd64/routestats

# compile 'windows'
env GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/routestats.exe
env GOOS=windows GOARCH=386 go build -o build/windows-386/routestats.exe
env GOOS=windows GOARCH=arm go build -o build/windows-arm/routestats.exe
