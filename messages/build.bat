protoc -I=. -I=%GOPATH%/src --gogoslick_out=plugins=grpc:.^
 list_trees.proto^
 create_tree.proto^
 delete_tree.proto^
 insert.proto^
 search.proto^
 remove.proto^
 traverse.proto^
 tree_identifier.proto^
 key_value_pair.proto
