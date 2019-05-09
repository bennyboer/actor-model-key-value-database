package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-sudo/messages"
	"testing"
	"time"
)

const timeout = time.Second * 5

// UTILITY METHODS

func getTreeServiceProps() *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {
		return newRoot()
	})
}

// Create a tree for testing.
func createTree(t *testing.T, rootContext *actor.RootContext, servicePID *actor.PID) *messages.TreeIdentifier {
	message := messages.CreateTreeRequest{}

	future := rootContext.RequestFuture(servicePID, &message, timeout)

	result, err := future.Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(*messages.CreateTreeResponse)
	if !ok {
		t.Errorf("expected response to be of type CreateTreeResponse")
	}

	if response.TreeId == nil {
		t.Errorf("expected non-nil tree identifier")
	}

	if response.TreeId.Id < 0 {
		t.Errorf("expected non-negative tree id")
	}

	if len(response.TreeId.Token) == 0 {
		t.Errorf("expected non-empty tree token")
	}

	return response.TreeId
}

func createTrees(t *testing.T, ctx *actor.RootContext, servicePID *actor.PID, count int) []*messages.TreeIdentifier {
	treeIds := make([]*messages.TreeIdentifier, 0, count)

	for i := 0; i < count; i++ {
		treeIds = append(treeIds, createTree(t, ctx, servicePID))
	}

	return treeIds
}

func listTrees(t *testing.T, ctx *actor.RootContext, servicePID *actor.PID) []*messages.TreeIdentifier {
	result, err := ctx.RequestFuture(servicePID, &messages.ListTreesRequest{}, timeout).Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(*messages.ListTreesResponse)
	if !ok {
		t.Errorf("expected response to be of type ListTreesResponse")
	}

	if response.TreeIds == nil {
		t.Errorf("expected response to have a slice of tree identifiers")
	}
	return response.TreeIds
}

func insert(t *testing.T, ctx *actor.RootContext, servicePID *actor.PID, treeId *messages.TreeIdentifier, key int, value string) *messages.InsertResponse {
	result, err := ctx.RequestFuture(servicePID, &messages.InsertRequest{
		TreeId: treeId,
		Entry: &messages.KeyValuePair{
			Key:   int32(key),
			Value: value,
		},
	}, timeout).Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(*messages.InsertResponse)
	if !ok {
		t.Errorf("expected result to be of type InsertResponse^")
	}

	if !response.Success {
		t.Errorf("expected insert to be successful")
	}

	return response
}

func search(t *testing.T, ctx *actor.RootContext, servicePID *actor.PID, treeId *messages.TreeIdentifier, key int) *messages.SearchResponse {
	result, err := ctx.RequestFuture(servicePID, &messages.SearchRequest{
		Key:    int32(key),
		TreeId: treeId,
	}, timeout).Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(*messages.SearchResponse)
	if !ok {
		t.Errorf("expected search result to be of type SearchResponse")
	}

	return response
}

func remove(t *testing.T, ctx *actor.RootContext, servicePID *actor.PID, treeId *messages.TreeIdentifier, key int) *messages.RemoveResponse {
	result, err := ctx.RequestFuture(servicePID, &messages.RemoveRequest{
		TreeId: treeId,
		Key:    int32(key),
	}, timeout).Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(*messages.RemoveResponse)
	if !ok {
		t.Errorf("expected result to be of type RemoveResponse")
	}

	return response
}

func traverse(t *testing.T, ctx *actor.RootContext, servicePID *actor.PID, treeId *messages.TreeIdentifier) []*messages.KeyValuePair {
	result, err := ctx.RequestFuture(servicePID, &messages.TraverseRequest{
		TreeId: treeId,
	}, timeout).Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(*messages.TraverseResponse)
	if !ok {
		t.Errorf("expected result to have type TraverseResponse")
	}

	return response.Pairs
}

// TEST METHODS

func TestService_CreateTree(t *testing.T) {
	rootContext := actor.EmptyRootContext
	servicePID := rootContext.Spawn(getTreeServiceProps())

	createTree(t, rootContext, servicePID)
	// Well, if it works it ain't stupid
}

func TestService_ListTrees(t *testing.T) {
	rootContext := actor.EmptyRootContext
	servicePID := rootContext.Spawn(getTreeServiceProps())

	// Create multiple trees first
	expectedTreeIds := createTrees(t, rootContext, servicePID, 3)
	actualTreeIds := listTrees(t, rootContext, servicePID)

	for i := 0; i < 3; i++ {
		expected := expectedTreeIds[i]
		actual := actualTreeIds[i]
		// TODO unterschiedliche Längen! Actor abschießen oder kompensieren

		if expected.Id != actual.Id {
			t.Errorf("expected tree id %d; got %d", expected.Id, actual.Id)
		}

		if len(actual.Token) != 0 {
			t.Errorf("tree token mustn't be revealed by the action! It is a secret!")
		}
	}

}

func TestService_DeleteTree(t *testing.T) {
	rootContext := actor.EmptyRootContext
	servicePID := rootContext.Spawn(getTreeServiceProps())

	// Create multiple trees first
	treeIds := createTrees(t, rootContext, servicePID, 3)

	// Delete one
	treeToDelete := treeIds[1]
	result, err := rootContext.RequestFuture(servicePID, &messages.DeleteTreeRequest{
		TreeId: treeToDelete,
	}, timeout).Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok := result.(*messages.DeleteTreeResponse)
	if !ok {
		t.Errorf("expected result to be of type DeleteTreeResponse")
	}

	if !response.Success {
		t.Errorf("expected response to be successful")
	}

	if !response.MarkedForDeletion {
		t.Errorf("expected tree to be marked for deletion")
	}

	// Check that tree has not yet been deleted
	remainingTreeIds := listTrees(t, rootContext, servicePID)
	if len(remainingTreeIds) != 3 {
		t.Errorf("expected all 3 trees to be still there")
	}

	// Try again to delete it forever
	result, err = rootContext.RequestFuture(servicePID, &messages.DeleteTreeRequest{
		TreeId: treeToDelete,
	}, timeout).Result()
	if err != nil {
		t.Errorf("expected no error")
	}

	response, ok = result.(*messages.DeleteTreeResponse)
	if !ok {
		t.Errorf("expected result to be of type DeleteTreeResponse")
	}

	if !response.Success {
		t.Errorf("expected deletion to be successful")
	}

	if response.MarkedForDeletion {
		t.Errorf("expected tree to be gone by now and not been marked for deletion, since it already was")
	}

	// Check that tree is really gone
	remainingTreeIds = listTrees(t, rootContext, servicePID)
	if len(remainingTreeIds) != 2 {
		t.Errorf("expected one tree to be gone by now")
	}

	for _, remainingTreeId := range remainingTreeIds {
		if remainingTreeId.Id == treeToDelete.Id {
			t.Errorf("this tree should have been deleted by now")
		}
	}
}

func TestService_InsertKeyValuePair(t *testing.T) {
	rootContext := actor.EmptyRootContext
	servicePID := rootContext.Spawn(getTreeServiceProps())

	// Create tree to test first
	treeId := createTree(t, rootContext, servicePID)

	entry := messages.KeyValuePair{
		Key:   54,
		Value: "Ich bin ein Wert! :)",
	}

	insert(t, rootContext, servicePID, treeId, int(entry.Key), entry.Value)

	// Check if entry has really been inserted
	searchResponse := search(t, rootContext, servicePID, treeId, int(entry.Key))
	if searchResponse.Success != true {
		t.Errorf("previously inserted key-value pair could not be found")
	}

	if searchResponse.Entry.Key != entry.Key {
		t.Errorf("the inserted entry key %d and the found key %d do not match", searchResponse.Entry.Key, entry.Key)
	}

	if searchResponse.Entry.Value != entry.Value {
		t.Errorf("the inserted entry value '%s' and the found value '%s' do not match", searchResponse.Entry.Value, entry.Value)
	}
}

func TestService_SearchKeyValuePair(t *testing.T) {
	rootContext := actor.EmptyRootContext
	servicePID := rootContext.Spawn(getTreeServiceProps())

	// Create tree to test first
	treeId := createTree(t, rootContext, servicePID)

	// Test if something, that is not there can't be found
	result := search(t, rootContext, servicePID, treeId, 54)

	if result.Success {
		t.Errorf("expected search to be unsuccessful")
	}

	// The rest has already been tested with the Insert action test.
}

func TestService_RemoveKeyValuePair(t *testing.T) {
	rootContext := actor.EmptyRootContext
	servicePID := rootContext.Spawn(getTreeServiceProps())

	// Create tree to test first
	treeId := createTree(t, rootContext, servicePID)

	// Insert value
	insert(t, rootContext, servicePID, treeId, 3243, "Hallo Welt")

	// Delete it
	result := remove(t, rootContext, servicePID, treeId, 3243)
	if !result.Success {
		t.Errorf("expected remove to be successful")
	}

	if result.RemovedPair == nil {
		t.Errorf("expected remove to return the removed entry")
	}

	if result.RemovedPair.Key != int32(3243) {
		t.Errorf("expected removed entry to have the correct key %d and not %d", int32(3243), result.RemovedPair.Key)
	}

	if result.RemovedPair.Value != "Hallo Welt" {
		t.Errorf("expected removed pair to have the correct value")
	}

	// Check if it can be still found
	searchResult := search(t, rootContext, servicePID, treeId, 3243)
	if searchResult.Success {
		t.Errorf("expected entry to have been removed")
	}
}

func TestService_TraverseKeyValuePairs(t *testing.T) {
	rootContext := actor.EmptyRootContext
	servicePID := rootContext.Spawn(getTreeServiceProps())

	// First and foremost create a tree
	treeId := createTree(t, rootContext, servicePID)

	// Then create some entries
	count := 13
	expectedEntries := make([]*messages.KeyValuePair, 0, count)
	for i := 0; i < count; i++ {
		expectedEntries = append(expectedEntries, &messages.KeyValuePair{
			Key:   int32(i + 1),
			Value: fmt.Sprintf("Value %d", i+1),
		})
	}

	// Insert entries into tree
	for _, entry := range expectedEntries {
		insert(t, rootContext, servicePID, treeId, int(entry.Key), entry.Value)
	}

	entries := traverse(t, rootContext, servicePID, treeId)
	if entries == nil {
		t.Errorf("expected some entires and not nil")
	}

	if len(entries) != count {
		t.Errorf("expected %d entries; got %d", count, len(entries))
	}

	for i := 0; i < count; i++ {
		expected := expectedEntries[i]
		actual := entries[i]

		if expected.Key != actual.Key {
			t.Errorf("expected key %d; got %d", expected.Key, actual.Key)
		}

		if expected.Value != actual.Value {
			t.Errorf("expected value '%s'; got '%s'", expected.Value, actual.Value)
		}
	}
}
