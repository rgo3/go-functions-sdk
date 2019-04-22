package auth

import (
	"context"
	"fmt"

	"github.com/dergoegge/go-functions-sdk/pkg/functions"
)

var AuthFunc = functions.New().Auth().OnCreate(func(ctx context.Context, event functions.AuthEvent) error {
	fmt.Printf("New user %s", event.Email)
	return nil
})
