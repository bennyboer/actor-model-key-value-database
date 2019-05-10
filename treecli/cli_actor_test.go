package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	act "github.com/ob-vss-ss19/blatt-3-sudo/treecli/action"
	"github.com/ob-vss-ss19/blatt-3-sudo/treecli/localmessages"
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
		&localmessages.CLIExecuteRequest{
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

	response, ok := result.(*localmessages.CLIExecuteReply)
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
		"5",
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
		&localmessages.CLIExecuteRequest{
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

	response, ok := result.(*localmessages.CLIExecuteReply)
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

func TestCLIActor_ExecuteState_DeleteTree_InvalidInput(t *testing.T) {
	arguments := []string{}[:]
	flags := util.Flags{
		ID: -1,
	}

	action := act.DeleteTree{}
	err := action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected to throw error because of negative tree id")
	}

	flags.ID = 123

	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected to throw error because of missing tree token")
	}
}

func TestCLIActor_ExecuteState_DeleteTree(t *testing.T) {
	arguments := []string{
		"delete-tree",
	}
	flags := util.Flags{
		Timeout: time.Second * 5,
		ID:      123,
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
		&localmessages.CLIExecuteRequest{
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

	response, ok := result.(*localmessages.CLIExecuteReply)
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
		&localmessages.CLIExecuteRequest{
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

	response, ok = result.(*localmessages.CLIExecuteReply)
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
	flags.ID = 324
	future = rootContext.RequestFuture(
		cliPID,
		&localmessages.CLIExecuteRequest{
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

	response, ok = result.(*localmessages.CLIExecuteReply)
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

func TestCLIActor_ExecuteState_Insert_ValidateInput(t *testing.T) {
	arguments := []string{}[:]
	flags := util.Flags{
		ID: -1,
	}

	// Test negative tree id
	action := act.Insert{}
	err := action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected action to throw an error because of a negative tree id")
	}

	// Test empty token
	flags.ID = 1
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected action to throw an error because of an empty token")
	}

	// Test invalid argument count
	flags.Token = "hw323"
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of insufficient length of the argument slice")
	}

	arguments = []string{
		"1",
	}[:]
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of insufficient length of the argument slice")
	}

	// Test non-integer key
	arguments = []string{
		"key",
		"value",
	}[:]
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of non-integer key")
	}

	// Test nil context
	arguments = []string{
		"1",
		"\"Hello",
		"World\"",
	}[:]

	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because context was nil")
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
		ID:      123,
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
		&localmessages.CLIExecuteRequest{
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

	response, ok := result.(*localmessages.CLIExecuteReply)
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
	flags.ID = 5435
	future = rootContext.RequestFuture(
		cliPID,
		&localmessages.CLIExecuteRequest{
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

	response, ok = result.(*localmessages.CLIExecuteReply)
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

func TestCLIActor_ExecuteState_Remove_ValidateInputs(t *testing.T) {
	arguments := []string{}[:]
	flags := util.Flags{
		ID: -1,
	}

	action := act.Remove{}

	// Test negative tree id
	err := action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of negative tree id")
	}

	// Test empty token
	flags.ID = 123
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of empty tree token")
	}

	// Test correct arguments slice length
	flags.Token = "abc123"
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of incorrect arguments length")
	}

	arguments = []string{
		"argument1",
		"argument2",
	}[:]
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of incorrect arguments length")
	}

	// Test non-integer key
	arguments = []string{
		"key",
	}[:]
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because 'key' is not an integer")
	}
}

func TestCLIActor_ExecuteState_Remove(t *testing.T) {
	arguments := []string{
		"remove",
		"434",
	}
	flags := util.Flags{
		ID:      123,
		Token:   "abc123",
		Timeout: time.Second * 5,
	}

	rootContext := actor.EmptyRootContext

	serviceProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.RemoveRequest:
			fmt.Printf("incoming message %v\n", msg)

			success := true
			if msg.Key != 1 {
				success = false
			}

			var removedPair *messages.KeyValuePair = nil
			if success {
				removedPair = &messages.KeyValuePair{
					Key:   1,
					Value: "Hallo Welt",
				}
			}

			ctx.Respond(&messages.RemoveResponse{
				Success:     success,
				RemovedPair: removedPair,
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
		&localmessages.CLIExecuteRequest{
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

	response, ok := result.(*localmessages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage := "Could not remove key-value pair."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok := response.Original.(*messages.RemoveResponse)
	if !ok {
		t.Errorf("expected original response to be of type RemoveResponse")
	}

	if original.Success != false {
		t.Errorf("expected original response to be unsuccessful")
	}

	if original.RemovedPair != nil {
		t.Errorf("expected original response to have nil pointer as removed pair")
	}

	// Test successful delete
	arguments = []string{
		"remove",
		"1",
	}

	future = rootContext.RequestFuture(
		cliPID,
		&localmessages.CLIExecuteRequest{
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

	response, ok = result.(*localmessages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage = "Removed key-value pair with Key: 1 and Value 'Hallo Welt'."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok = response.Original.(*messages.RemoveResponse)
	if !ok {
		t.Errorf("expected original response to be of type RemoveResponse")
	}

	if original.Success != true {
		t.Errorf("expected original response to be successful")
	}

	if original.RemovedPair == nil || original.RemovedPair.Key != 1 || original.RemovedPair.Value != "Hallo Welt" {
		t.Errorf("expected original response to have a removed pair with key: 1 and value: 'Hallo Welt'")
	}
}

func TestCLIActor_ExecuteState_Search(t *testing.T) {
	arguments := []string{
		"search",
		"345",
	}
	flags := util.Flags{
		ID:      123,
		Token:   "abc123",
		Timeout: time.Second * 5,
	}

	rootContext := actor.EmptyRootContext

	serviceProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.SearchRequest:
			fmt.Printf("incoming message %v\n", msg)

			success := true
			if msg.Key != 1 {
				success = false
			}

			var pair *messages.KeyValuePair = nil
			if success {
				pair = &messages.KeyValuePair{
					Key:   1,
					Value: "Hallo Welt",
				}
			}

			ctx.Respond(&messages.SearchResponse{
				Success: success,
				Entry:   pair,
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
		&localmessages.CLIExecuteRequest{
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

	response, ok := result.(*localmessages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage := "Could not find key-value pair."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok := response.Original.(*messages.SearchResponse)
	if !ok {
		t.Errorf("expected original response to be of type SearchResponse")
	}

	if original.Success != false {
		t.Errorf("expected original response to be unsuccessful")
	}

	if original.Entry != nil {
		t.Errorf("expected original response to have nil pointer as found pair")
	}

	// Test successful search
	arguments = []string{
		"search",
		"1",
	}

	future = rootContext.RequestFuture(
		cliPID,
		&localmessages.CLIExecuteRequest{
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

	response, ok = result.(*localmessages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage = "Found key-value pair with Key: 1 and Value 'Hallo Welt'."

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok = response.Original.(*messages.SearchResponse)
	if !ok {
		t.Errorf("expected original response to be of type SearchResponse")
	}

	if original.Success != true {
		t.Errorf("expected original response to be successful")
	}

	if original.Entry == nil || original.Entry.Key != 1 || original.Entry.Value != "Hallo Welt" {
		t.Errorf("expected original response to have a found pair with key: 1 and value: 'Hallo Welt'")
	}
}

func TestCLIActor_ExecuteState_Search_ValidateInputs(t *testing.T) {
	arguments := []string{}[:]
	flags := util.Flags{
		ID: -1,
	}

	action := act.Search{}

	// Test negative tree id
	err := action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of negative tree id")
	}

	// Test empty token
	flags.ID = 123
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of empty tree token")
	}

	// Test correct arguments slice length
	flags.Token = "abc123"
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of incorrect arguments length")
	}

	arguments = []string{
		"argument1",
		"argument2",
	}[:]
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of incorrect arguments length")
	}

	// Test non-integer key
	arguments = []string{
		"key",
	}[:]
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because 'key' is not an integer")
	}
}

func TestCLIActor_ExecuteState_Traverse(t *testing.T) {
	arguments := []string{
		"traverse",
	}
	flags := util.Flags{
		ID:      124,
		Token:   "abc123",
		Timeout: time.Second * 5,
	}

	rootContext := actor.EmptyRootContext

	serviceProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.TraverseRequest:
			fmt.Printf("incoming message %v\n", msg)

			if msg.TreeId.Id == 123 {
				ctx.Respond(&messages.TraverseResponse{
					Success: true,
					Pairs: []*messages.KeyValuePair{
						&messages.KeyValuePair{
							Key:   1,
							Value: "Key 1",
						},
						&messages.KeyValuePair{
							Key:   2,
							Value: "Key 2",
						},
						&messages.KeyValuePair{
							Key:   3,
							Value: "Key 3",
						},
					}[:],
				})
			} else {
				ctx.Respond(&messages.TraverseResponse{
					Pairs:        nil,
					Success:      false,
					ErrorMessage: "Tree id is incorrect",
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
		&localmessages.CLIExecuteRequest{
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

	response, ok := result.(*localmessages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage := "Tree id is incorrect"

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok := response.Original.(*messages.TraverseResponse)
	if !ok {
		t.Errorf("expected original response to be of type TraverseResponse")
	}

	if original.Pairs != nil {
		t.Errorf("expected original response to have no pairs")
	}

	if original.Success != false {
		t.Error("expected original response to be unsuccessful")
	}

	// Test successful traverse
	flags.ID = 123

	future = rootContext.RequestFuture(
		cliPID,
		&localmessages.CLIExecuteRequest{
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

	response, ok = result.(*localmessages.CLIExecuteReply)
	if !ok {
		t.Errorf("expected response to be of type CLIExecuteReply")
	}

	expectedResultMessage = `All key-value pairs:
  - Key: 1, Value: 'Key 1'
  - Key: 2, Value: 'Key 2'
  - Key: 3, Value: 'Key 3'
`

	if response.Message != expectedResultMessage {
		t.Errorf("expected message '%s'; got '%s'", expectedResultMessage, response.Message)
	}

	original, ok = response.Original.(*messages.TraverseResponse)
	if !ok {
		t.Errorf("expected original response to be of type TraverseResponse")
	}

	if original.Pairs == nil {
		t.Errorf("expected original response to have pairs")
	}

	if original.Success != true {
		t.Error("expected original response to be successful")
	}

	if len(original.Pairs) != 3 {
		t.Errorf("expected original response to have 3 pairs as result")
	}
}

func TestCLIActor_ExecuteState_Traverse_ValidateInputs(t *testing.T) {
	arguments := []string{}[:]
	flags := util.Flags{
		ID: -1,
	}

	action := act.Traverse{}

	// Test negative tree id
	err := action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of negative tree id")
	}

	// Test empty token
	flags.ID = 123
	err = action.Execute(nil, &flags, arguments, nil)
	if err == nil {
		t.Errorf("expected error because of empty tree token")
	}
}
