package functions

import (
	"context"
	"time"
)

// Firestore event names
const (
	FirestoreOnCreate string = "providers/cloud.firestore/eventTypes/document.create"
	FirestoreOnUpdate string = "providers/cloud.firestore/eventTypes/document.update"
	FirestoreOnDelete string = "providers/cloud.firestore/eventTypes/document.delete"
	FirestoreOnWrite  string = "providers/cloud.firestore/eventTypes/document.write"
)

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time   `json:"createTime"`
	Fields     interface{} `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

// FirestoreFunctionBuilder holds a FunctionBuilder and the firestore event entrypoint
type FirestoreFunctionBuilder struct {
	FunctionBuilder *FunctionBuilder

	Doc     string
	Event   string
	Handler func(ctx context.Context, event FirestoreEvent) error
}

// FirestoreFunction is the function signature of a firestore event triggered function
type FirestoreFunction func(ctx context.Context, event FirestoreEvent) error

// Document sets the document for the function
// https://cloud.google.com/functions/docs/calling/cloud-firestore#specifying_the_document_path
func (fb *FirestoreFunctionBuilder) Document(doc string) *FirestoreFunctionBuilder {
	fb.Doc = doc
	return fb
}

// OnCreate creates a firestore function triggered on document creation
func (fb *FirestoreFunctionBuilder) OnCreate(fn FirestoreFunction) *FirestoreFunctionBuilder {
	return fb.onOperation(FirestoreOnCreate, fn)
}

// OnUpdate creates a firestore function triggered on a document update
func (fb *FirestoreFunctionBuilder) OnUpdate(fn FirestoreFunction) *FirestoreFunctionBuilder {
	return fb.onOperation(FirestoreOnUpdate, fn)
}

// OnDelete creates a firestore function triggered on document deletion
func (fb *FirestoreFunctionBuilder) OnDelete(fn FirestoreFunction) *FirestoreFunctionBuilder {
	return fb.onOperation(FirestoreOnDelete, fn)
}

// OnWrite creates a firestore function triggered on document write (creation, update, delete)
func (fb *FirestoreFunctionBuilder) OnWrite(fn FirestoreFunction) *FirestoreFunctionBuilder {
	return fb.onOperation(FirestoreOnWrite, fn)
}

func (fb *FirestoreFunctionBuilder) onOperation(event string, fn FirestoreFunction) *FirestoreFunctionBuilder {
	fb.Handler = fn
	fb.Event = event
	return fb
}
