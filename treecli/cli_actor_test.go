package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/local_messages"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/util"
	"testing"
	"time"
)

func TestNewCLIActor(t *testing.T) {
	a := NewCLIActor()

	if a == nil {
		t.Errorf("expected actor instance, not nil")
	}
}

func TestCLIActor_ExecuteState_ListTrees(t *testing.T) {
	arguments := []string{
		"list",
	}
	flags := util.Flags{
		Timeout: time.Second * 5,
	}

	rootContext := actor.EmptyRootContext

	serviceProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.ListTreesRequest:
			fmt.Printf("incoming message %v\n", msg)

			ctx.Respond(&messages.ListTreesResponse{
				TreeIds: []*messages.TreeIdentifier{
					&messages.TreeIdentifier{
						Id:    123,
						Token: "",
					},
					&messages.TreeIdentifier{
						Id:    234,
						Token: "",
					},
				}[:],
			})
		}
	})

	servicePID := rootContext.Spawn(serviceProps)

	cliProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCLIActor()
	})

	cliPID := rootContext.Spawn(cliProps)

	future := rootContext.RequestFuture(
		cliPID,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     &flags,
			RemotePID: servicePID,
		},
		flags.Timeout,
	)

	result, err := future.Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(local_messages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage := `Available Tree IDs:
  - 123
  - 234
`

	if response.Message != expectedResultMessage {
		t.Errorf("Expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok := response.Original.(*messages.ListTreesResponse)
	if !ok {
		t.Errorf("expected original response to be of type ListTreesResponse")
	}

	if len(original.TreeIds) != 2 {
		t.Errorf("expected original response to have two tree identifiers; got %d\n", len(original.TreeIds))
	}

	if original.TreeIds[0].Id != 123 || original.TreeIds[1].Id != 234 {
		t.Errorf("expected response does not have the correct content")
	}

	if original.TreeIds[0].Token != "" || original.TreeIds[1].Token != "" {
		t.Errorf("expected original response to not include Tokens, since they are secret!")
	}
}

func TestCLIActor_ExecuteState_CreateTree(t *testing.T) {
	arguments := []string{
		"create-tree",
	}
	flags := util.Flags{
		Timeout: time.Second * 5,
	}

	rootContext := actor.EmptyRootContext

	serviceProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.CreateTreeRequest:
			fmt.Printf("incoming message %v\n", msg)

			ctx.Respond(&messages.CreateTreeResponse{
				TreeId: &messages.TreeIdentifier{
					Id:    123,
					Token: "abc123",
				},
			})
		}
	})

	servicePID := rootContext.Spawn(serviceProps)

	cliProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCLIActor()
	})

	cliPID := rootContext.Spawn(cliProps)

	future := rootContext.RequestFuture(
		cliPID,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     &flags,
			RemotePID: servicePID,
		},
		flags.Timeout,
	)

	result, err := future.Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(local_messages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage := `Created tree with ID: 123 and Token: 'abc123'
`

	if response.Message != expectedResultMessage {
		t.Errorf("Expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok := response.Original.(*messages.CreateTreeResponse)
	if !ok {
		t.Errorf("expected original response to be of type CreateTreeResponse")
	}

	if original.TreeId.Id != 123 || original.TreeId.Token != "abc123" {
		t.Errorf("Expected response with tree id 123 and Token 'abc123'; got ID: %d and Token: %s", original.TreeId.Id, original.TreeId.Token)
	}
}

func TestCLIActor_ExecuteState_DeleteTree(t *testing.T) {
	arguments := []string{
		"delete-tree",
	}
	flags := util.Flags{
		Timeout: time.Second * 5,
		Id:      123,
		Token:   "abc123",
	}

	rootContext := actor.EmptyRootContext

	markedForDeletion := false

	serviceProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.DeleteTreeRequest:
			fmt.Printf("incoming message %v\n", msg)

			if msg.TreeId.Id != 123 {
				fmt.Println("Invalid tree id!")

				ctx.Respond(&messages.DeleteTreeResponse{
					Success:           false,
					MarkedForDeletion: false,
				})
			} else if markedForDeletion {
				fmt.Println("Delete forever!")

				ctx.Respond(&messages.DeleteTreeResponse{
					Success:           true,
					MarkedForDeletion: false,
				})
			} else {
				markedForDeletion = true

				fmt.Println("Marked for deletion!")

				ctx.Respond(&messages.DeleteTreeResponse{
					Success:           true,
					MarkedForDeletion: true,
				})
			}
		}
	})

	servicePID := rootContext.Spawn(serviceProps)

	cliProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCLIActor()
	})

	cliPID := rootContext.Spawn(cliProps)

	future := rootContext.RequestFuture(
		cliPID,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     &flags,
			RemotePID: servicePID,
		},
		flags.Timeout,
	)

	result, err := future.Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(local_messages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage := "Marked tree for deletion. Execute command again to delete tree forever."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok := response.Original.(*messages.DeleteTreeResponse)
	if !ok {
		t.Errorf("expected original response to be of type DeleteTreeResponse")
	}

	if original.MarkedForDeletion != true || original.Success != true {
		t.Errorf("expected MarkedToDeletion to be true and Success to be true as well")
	}

	// Try second deletion call to delete the tree forever.
	future = rootContext.RequestFuture(
		cliPID,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     &flags,
			RemotePID: servicePID,
		},
		flags.Timeout,
	)

	result, err = future.Result()
	if err != nil {
		t.Errorf("expected no error; got %s\n", err.Error())
	}

	response, ok = result.(local_messages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage = "Tree has been deleted."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok = response.Original.(*messages.DeleteTreeResponse)
	if !ok {
		t.Errorf("expected original response to be of type DeleteTreeResponse")
	}

	if original.MarkedForDeletion != false && original.Success != true {
		t.Errorf("expected the response to have MarkedForDeletion false and Success true")
	}

	// Try again, but now with an invalid tree id
	flags.Id = 324
	future = rootContext.RequestFuture(
		cliPID,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     &flags,
			RemotePID: servicePID,
		},
		flags.Timeout,
	)

	result, err = future.Result()
	if err != nil {
		t.Errorf("expected no error; got %s\n", err.Error())
	}

	response, ok = result.(local_messages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage = "Tree could not be deleted."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok = response.Original.(*messages.DeleteTreeResponse)
	if !ok {
		t.Errorf("expected original response to be of type DeleteTreeResponse")
	}

	if original.Success != false {
		t.Errorf("expected original response Success to be false")
	}
}

func TestCLIActor_ExecuteState_Insert(t *testing.T) {
	arguments := []string{
		"insert",
		"1",
		"Hello",
		"World",
	}
	flags := util.Flags{
		Timeout: time.Second * 5,
		Id:      123,
		Token:   "abc123",
	}

	rootContext := actor.EmptyRootContext

	serviceProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.InsertRequest:
			fmt.Printf("incoming message %v\n", msg)

			if msg.TreeId.Id == 123 {
				if msg.Entry.Key != 1 || msg.Entry.Value != "Hello World" {
					t.Errorf("expected key value pair to be {1, 'Hello World'}; got {%d, '%s'}", msg.Entry.Key, msg.Entry.Value)
				}

				ctx.Respond(&messages.InsertResponse{
					Success: true,
				})
			} else {
				ctx.Respond(&messages.InsertResponse{
					Success: false,
				})
			}
		}
	})

	servicePID := rootContext.Spawn(serviceProps)

	cliProps := actor.PropsFromProducer(func() actor.Actor {
		return NewCLIActor()
	})

	cliPID := rootContext.Spawn(cliProps)

	future := rootContext.RequestFuture(
		cliPID,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     &flags,
			RemotePID: servicePID,
		},
		flags.Timeout,
	)

	result, err := future.Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(local_messages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage := "Inserted successfully."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok := response.Original.(*messages.InsertResponse)
	if !ok {
		t.Errorf("expected original response to be of type InsertResponse")
	}

	if original.Success != true {
		t.Errorf("expected original response to be successful")
	}

	// Try again, but if invalid tree id
	flags.Id = 5435
	future = rootContext.RequestFuture(
		cliPID,
		&local_messages.CLIExecuteRequest{
			Arguments: arguments,
			Flags:     &flags,
			RemotePID: servicePID,
		},
		flags.Timeout,
	)

	result, err = future.Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok = result.(local_messages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage = "Insert did not work."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok = response.Original.(*messages.InsertResponse)
	if !ok {
		t.Errorf("expected original response to be of type InsertResponse")
	}

	if original.Success != false {
		t.Errorf("expected original response to have failed")
	}
}


