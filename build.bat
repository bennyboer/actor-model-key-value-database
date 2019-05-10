@echo off

echo Fetching dependencies
@echo on
go install github.com/gogo/protobuf/protoc-gen-gogoslick
go get
@echo off

echo Compiling protobuf messages
cd ./messages
@echo on
call ./build.bat
@echo off
cd ..

echo Building binaries
@echo on
go build -o bin/tree-cli.exe ./treecli
go build -o bin/tree-service.exe ./treeservice
@echo off
