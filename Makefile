docker: proto docker-treeservice docker-treecli

proto:
	cd messages && make regenerate-docker

docker-treeservice:
	docker build -f treeservice.dockerfile -t tree-service .

docker-treecli:
	docker build -f treecli.dockerfile -t tree-cli .
