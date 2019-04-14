@echo off

echo Fetching dependencies
@echo on
go get -t
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
