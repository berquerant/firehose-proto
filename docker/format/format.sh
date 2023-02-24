#!/bin/bash

if [ -n "$DRY_RUN" ] ; then
    # check
    opt="--dry-run --Werror"
else
    # overwrite
    opt="-i"
fi

find "$WORKDIR" -name "*.proto" -type f | grep -v firehose-docker-protobuf | xargs clang-format $opt --verbose
