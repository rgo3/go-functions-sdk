package functions

import (
	"context"
	"log"
	"time"
)

// Storage event names
const (
	StorageOnArchive        string = "google.storage.object.archive"
	StorageOnDelete         string = "google.storage.object.delete"
	StorageOnFinalize       string = "google.storage.object.finalize"
	StorageOnMetadataUpdate string = "google.storage.object.metadataUpdate"
)

// StorageEvent is the payload of a GCS event. Please refer to the docs for
// additional information regarding GCS events.
type StorageEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}

type StorageFunctionBuilder struct {
	FunctionBuilder *FunctionBuilder
	objSet          bool

	GCBucket string
	Event    string
	Handler  func(ctx context.Context, event StorageEvent) error
}

// StorageFunction holds the function signature of storage event triggered function
type StorageFunction func(ctx context.Context, event StorageEvent) error

// Bucket sets the trigger bucket
func (fb *StorageFunctionBuilder) Bucket(bucket string) *StorageFunctionBuilder {
	fb.GCBucket = bucket
	return fb
}

// Object is a no op but without it the function definition would read something like
// Bucket("bucket").OnDelete() which sounds like a trigger for the event of bucket deletion
// when really it is for the deletion of a object in the bucket
func (fb *StorageFunctionBuilder) Object() *StorageFunctionBuilder {
	fb.objSet = true
	return fb
}

// OnArchive is triggered when a live object is archived or deleted.
// Only triggered for versioning buckets.
// https://cloud.google.com/functions/docs/calling/storage#object_archive
func (fb *StorageFunctionBuilder) OnArchive(fn StorageFunction) *StorageFunctionBuilder {
	return fb.onOperation(StorageOnArchive, fn)
}

// OnDelete is triggered when a live object is deleted.
// For versioning buckets, this is only sent when a version is permanently deleted (but not when an object is archived).
// For non-versioning buckets, this is sent when an object is deleted or overwritten.
// https://cloud.google.com/functions/docs/calling/storage#object_delete
func (fb *StorageFunctionBuilder) OnDelete(fn StorageFunction) *StorageFunctionBuilder {
	return fb.onOperation(StorageOnDelete, fn)
}

// OnFinalize is triggered when a new object is created (or an existing object is overwritten,
// and a new generation of that object is created) in the bucket.
// https://cloud.google.com/functions/docs/calling/storage#object_finalize
func (fb *StorageFunctionBuilder) OnFinalize(fn StorageFunction) *StorageFunctionBuilder {
	return fb.onOperation(StorageOnFinalize, fn)
}

// OnMetadataUpdate is triggered when an existing object's metadata changes
// https://cloud.google.com/functions/docs/calling/storage#object_metadata_update
func (fb *StorageFunctionBuilder) OnMetadataUpdate(fn StorageFunction) *StorageFunctionBuilder {
	return fb.onOperation(StorageOnMetadataUpdate, fn)
}

func (fb *StorageFunctionBuilder) onOperation(event string, fn StorageFunction) *StorageFunctionBuilder {
	if !fb.objSet {
		log.Fatal("Object() was not called on the FunctionBuilder")
	}

	fb.Event = event
	fb.Handler = fn
	return fb
}
