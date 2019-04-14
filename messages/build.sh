#!/usr/bin/env bash
protoc -I=. -I=${GOPATH}/src --gogoslick_out=plugins=grpc:. \
    delete_tree.proto \
    insert.proto \
    key_value_pair.proto \
    remove.proto \
    search.proto \
    tree.proto \
    tree_identifier.proto \
    tree_service.proto \
    void.proto
