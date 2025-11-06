#!/usr/bin/env bash

set -e

# @cmd Run the binary without building
# @meta default-subcommand
run() {
    go run .
}

# @cmd Prepares a development environment
dev() {
    go mod tidy
    ls *.go | entr -c -r go run .
}

# @cmd Buils the binary
build() {
    go build
}

# See more details at https://github.com/sigoden/argc
eval "$(argc --argc-eval "$0" "$@")"
