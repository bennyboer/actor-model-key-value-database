#!/bin/sh

PROTOBUF_VERSION=3.7.1
PROTOC_FILENAME=protoc-${PROTOBUF_VERSION}-linux-x86_64.zip
wget https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/${PROTOC_FILENAME}
unzip ${PROTOC_FILENAME}
bin/protoc --version
ls bin

# Add to path
echo $PATH
PROTOC_BINARY_PATH=$(realpath bin/protoc)
export PATH=${PATH}:${PROTOC_BINARY_PATH}
echo $PATH
