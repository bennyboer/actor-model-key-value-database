#!/bin/sh

PROTOBUF_VERSION=3.7.1
PROTOC_FILENAME=protoc-${PROTOBUF_VERSION}-linux-x86_64.zip
wget -nc https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/${PROTOC_FILENAME}
unzip -o ${PROTOC_FILENAME}
chmod +x ./bin/protoc
./bin/protoc --version
ls bin

# Add to path
PROTOC_BINARY_PATH=$(realpath bin/protoc)
PATH=${PATH}:${PROTOC_BINARY_PATH}
echo $PATH
