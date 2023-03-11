#!/bin/bash

set -e

thisd="$(cd $(dirname $0);pwd)"
rootd="${thisd}/.."
cd "$rootd"

git ls-files | grep -E ".pb.go$" | grep -v grpc | while read target ; do
    echo "gen-proto: ${target}"
    make "$target"
done
