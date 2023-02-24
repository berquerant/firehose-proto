#!/bin/bash

thisd="$(cd $(dirname $0);pwd)"
rootd="${thisd}/.."
internald="/usr/app/src"

docker run -v "${rootd}:/${internald}" -w "${internald}" --rm yoheimuta/protolint lint .
