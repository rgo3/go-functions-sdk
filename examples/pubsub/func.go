package pubsub

import (
	"context"
	"fmt"

	"github.com/dergoegge/go-functions-sdk/pkg/functions"
)

var PubSubFunc = functions.New().PubSub().Topic("test-topic").
	OnPublish(func(ctx context.Context, event functions.PubSubMessage) error {
		fmt.Println(string(event.Data))
		return nil
	})
