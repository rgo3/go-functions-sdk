package functions

import (
	"encoding/json"
	"net/http"
)

// RuntimeOptions holds a timeout and memory option
type RuntimeOptions struct {
	Timeout int // number of seconds
	Memory  string
}

// FunctionBuilder holds the deployment config for a function
type FunctionBuilder struct {
	Retry       bool
	GCRegion    string
	RuntimeOpts RuntimeOptions
}

// HTTPFunctionBuilder holds a FunctionBuilder and the http entry point
type HTTPFunctionBuilder struct {
	FunctionBuilder *FunctionBuilder

	Handler func(w http.ResponseWriter, r *http.Request)
}

// New creates a new default FunctionBuilder
func New() *FunctionBuilder {
	return &FunctionBuilder{
		Retry:    false,
		GCRegion: "us-central1",
		RuntimeOpts: RuntimeOptions{
			Timeout: 60,
			Memory:  "256MB",
		},
	}
}

func (fb *FunctionBuilder) String() string {
	json, _ := json.Marshal(fb)
	return string(json)
}

// Region sets a region for a FunctionBuilder
func (fb *FunctionBuilder) Region(region string) *FunctionBuilder {
	fb.GCRegion = region
	return fb
}

// RetryOnFailure sets the retry option for deployment
func (fb *FunctionBuilder) RetryOnFailure(retry bool) *FunctionBuilder {
	fb.Retry = retry
	return fb
}

// RunWith sets runtime options for the function
func (fb *FunctionBuilder) RunWith(opts RuntimeOptions) *FunctionBuilder {
	fb.RuntimeOpts = opts
	return fb
}

// OnRequest creates a new http triggered function
func (fb *FunctionBuilder) OnRequest(fn func(w http.ResponseWriter, r *http.Request)) *HTTPFunctionBuilder {
	return &HTTPFunctionBuilder{
		FunctionBuilder: fb,
		Handler:         fn,
	}
}

// Firestore creates a FirestoreFunctionBuilder
func (fb *FunctionBuilder) Firestore() *FirestoreFunctionBuilder {
	return &FirestoreFunctionBuilder{
		FunctionBuilder: fb,
	}
}

// Storage creates a StorageFunctionBuilder
func (fb *FunctionBuilder) Storage() *StorageFunctionBuilder {
	return &StorageFunctionBuilder{
		FunctionBuilder: fb,
	}
}
