package firestore

import (
	"context"
	"fmt"

	"github.com/dergoegge/go-functions-sdk/pkg/functions"
)

// FirestoreFunc is a firstore triggered function
var FirestoreFunc = functions.New().
	// Set runtime options for the function (optional)
	// but if runtime options are set both timeout and memory must be set.
	RunWith(functions.RuntimeOptions{
		// Function will time out after 70 seconds
		Timeout: 70,
		// Function will run with 128MB of memory
		Memory: "128MB",
	}).
	// Set the region to europes-west2 (optional)
	Region("europe-west2").
	// Set firestore trigger to onwrite events on documents in the "test-collection"
	Firestore().Document("test-collection/{docID}").
	OnWrite(func(ctx context.Context, event functions.FirestoreEvent) error {
		fmt.Println(event.Value.Fields)
		return nil
	})
