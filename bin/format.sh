#!/bin/bash

thisd="$(cd $(dirname $0);pwd)"
rootd="${thisd}/.."
internald="/usr/app/src"
docker run -v "${rootd}:${internald}" -w "${internald}" -e WORKDIR="${internald}" -e DRY_RUN="$1" --rm firehose-proto-format format.sh
