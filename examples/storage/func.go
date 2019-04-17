package storage

import (
	"context"
	"fmt"

	"github.com/dergoegge/go-functions-sdk/pkg/functions"
)

// StorageFunc is a storage bucket triggered function that uses the default region and runtime options
// Triggered on the default bucket when a new object is created or an object is updated.
var StorageFunc = functions.New().
	Storage().Object().
	OnFinalize(func(ctx context.Context, event functions.StorageEvent) error {
		fmt.Println(event.Name)
		return nil
	})
