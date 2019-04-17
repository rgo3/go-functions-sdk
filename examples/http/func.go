package http

import (
	"fmt"
	"net/http"

	"github.com/dergoegge/go-functions-sdk/pkg/functions"
)

// HTTPFunc is a http triggered function that uses the default region and runtime options.
var HTTPFunc = functions.New().OnRequest(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello http func!")
})
