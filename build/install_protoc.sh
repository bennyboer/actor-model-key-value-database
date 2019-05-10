#!/bin/sh

PROTOBUF_VERSION=3.7.1
PROTOC_FILENAME=protoc-${PROTOBUF_VERSION}-linux-x86_64.zip
if [[ ! -e ${PROTOC_FILENAME} ]]
then
    wget https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/${PROTOC_FILENAME}
    unzip -o ${PROTOC_FILENAME}
fi

chmod +x ./bin/protoc

# Add to path
PROTOC_BINARY_PATH=$(realpath bin/protoc)
PATH=${PATH}:${PROTOC_BINARY_PATH}
echo $PATH

protoc --version
