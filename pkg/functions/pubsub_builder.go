package functions

import "context"

// PubSub event names
const (
	PubSubPublish = "google.pubsub.topic.publish"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type PubSubFunctionBuilder struct {
	FunctionBuilder *FunctionBuilder

	Event   string
	PSTopic string
	Handler func(ctx context.Context, event PubSubMessage) error
}

// PubSubFunction holds the function signature of a pubsub event
type PubSubFunction func(ctx context.Context, event PubSubMessage) error

// Topic sets the trigger topic
func (fb *PubSubFunctionBuilder) Topic(topic string) *PubSubFunctionBuilder {
	fb.PSTopic = topic
	return fb
}

// OnPublish is triggered when a message is published to the Cloud Pub/Sub topic that is specified when a function is deployed.
// Every message published to this topic will trigger function execution with message contents passed as input data.
func (fb *PubSubFunctionBuilder) OnPublish(fn PubSubFunction) *PubSubFunctionBuilder {
	fb.Event = PubSubPublish
	fb.Handler = fn
	return fb
}
