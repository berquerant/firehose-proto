FROM debian:11.6-slim

COPY ./format.sh /usr/local/bin/format.sh

RUN apt-get update && apt-get install -y \
    clang-format \
 && rm -rf /var/lib/apt/lists/* && \
    clang-format --version
