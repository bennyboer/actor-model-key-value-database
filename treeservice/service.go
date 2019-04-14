package main

import (
	"context"
	"errors"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"log"
)

// The service implementation
type Service struct{}

func (s *Service) ListTrees(in *messages.VoidMessage, stream messages.TreeService_ListTreesServer) error {
	log.Println("ListTrees called")

	test := []int32{723, 256, 752, 647, 645, 567}

	for _, id := range test {
		if err := stream.Send(&messages.TreeIdentifier{Id: id}); err != nil {
			return err
		}
	}

	// TODO Implement the real functionality

	return nil
}

func (Service) CreateTree(ctx context.Context, in *messages.VoidMessage) (*messages.TreeIdentifier, error) {
	log.Println("CreateTree called")

	return nil, errors.New("not yet implemented :(")
}

func (Service) DeleteTree(ctx context.Context, in *messages.TreeIdentifier) (*messages.DeleteTreeReply, error) {
	log.Println("DeleteTree called")

	return nil, errors.New("not yet implemented :(")
}

func (Service) Insert(ctx context.Context, in *messages.InsertRequest) (*messages.VoidMessage, error) {
	log.Println("Insert called")

	return nil, errors.New("not yet implemented :(")
}

func (Service) Remove(ctx context.Context, in *messages.RemoveRequest) (*messages.VoidMessage, error) {
	log.Println("Remove called")

	return nil, errors.New("not yet implemented :(")
}

func (Service) Search(ctx context.Context, in *messages.SearchRequest) (*messages.SearchResponse, error) {
	log.Println("Search called")

	return nil, errors.New("not yet implemented :(")
}

func (Service) Traverse(treeId *messages.TreeIdentifier, stream messages.TreeService_TraverseServer) error {
	log.Println("Traverse called")

	return errors.New("not yet implemented :(")
}
