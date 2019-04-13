#!/usr/bin/env bash

echo "Fetching dependencies"
go get -t

echo "Compiling protobuf messages"
cd ./messages
./build.sh
cd ..

echo "Building binaries"
go build -o bin/tree-cli.exe ./treecli/main.go
go build -o bin/tree-service.exe ./treeservice/main.go
