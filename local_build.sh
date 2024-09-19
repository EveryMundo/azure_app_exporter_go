#!/usr/bin/env sh

set -eu

docker build -t go-builder .

docker run -it --rm -v ./:/src go-builder go build
