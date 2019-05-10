regenerate-docker:
	docker run --rm -it \
		-v `pwd`:/messages \
		-w /messages \
		obraun/vss-protoactor-jenkins \
		make regenerate

regenerate:
	go install -v github.com/gogo/protobuf/protoc-gen-gogoslick
	protoc -I=. --gogoslick_out=plugins=grpc:. \
    	list_trees.proto \
    	create_tree.proto \
    	delete_tree.proto \
    	insert.proto \
    	search.proto \
    	remove.proto \
    	traverse.proto \
    	tree_identifier.proto \
    	key_value_pair.proto
