regenerate-docker:
	docker run --rm -it \
		-v `pwd`:/messages \
		-w /messages \
		obraun/vss-protoactor-jenkins \
		make regenerate

regenerate:
	go install -v github.com/gogo/protobuf/protoc-gen-gogoslick
	protoc -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
		--gogoslick_out=plugins=grpc:. \
		delete_tree.proto \
        insert.proto \
        key_value_pair.proto \
        remove.proto \
        search.proto \
        tree.proto \
        tree_identifier.proto \
        tree_service.proto \
        void.proto
