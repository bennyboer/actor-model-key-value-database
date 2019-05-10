#!/bin/sh

echo "Fetching dependencies"
go clean
go install github.com/gogo/protobuf/protoc-gen-gogoslick
go get

echo "Compiling protobuf messages"
cd ./messages
chmod +x ./build.sh
./build.sh
cd ..

echo "Building binaries"
go build -o bin/tree-cli.exe ./treecli
go build -o bin/tree-service.exe ./treeservice
