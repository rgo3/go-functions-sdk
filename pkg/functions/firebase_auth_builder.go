package functions

import (
	"context"
	"time"
)

// Firebase authentication event names
const (
	AuthUserCreate string = "providers/firebase.auth/eventTypes/user.create"
	AuthUserDelete string = "providers/firebase.auth/eventTypes/user.delete"
)

type FirebaseAuthFunctionBuilder struct {
	FunctionBuilder *FunctionBuilder

	Event   string
	Handler func(ctx context.Context, event AuthEvent) error
}

// FirebaseAuthFunction holds the function signature of a firebase auth triggered function
type FirebaseAuthFunction func(ctx context.Context, event AuthEvent) error

// AuthEvent is the payload of a Firestore Auth event.
type AuthEvent struct {
	UID string `json:"uid"`

	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	DisplayName   string `json:"displayName"`
	PhotoURL      string `json:"photoURL"`

	ProviderData struct {
		UID         string `json:"uid"`
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
		PhotoURL    string `json:"photoURL"`
		ProviderID  string `json:"providerId"`
	} `json:"providerData"`

	CustomClaims *interface{} `json:"customClaims"`

	Disabled bool `json:"disabled"`
	Metadata struct {
		CreatedAt      time.Time `json:"createdAt"`
		LastSignInTime time.Time `json:"lastSignInTime"`
	} `json:"metadata"`
}

// OnCreate is triggered when a user account is created.
func (fb *FirebaseAuthFunctionBuilder) OnCreate(fn FirebaseAuthFunction) *FirebaseAuthFunctionBuilder {
	return fb.onOperation(AuthUserCreate, fn)
}

// OnDelete is triggered when a user account is deleted.
func (fb *FirebaseAuthFunctionBuilder) OnDelete(fn FirebaseAuthFunction) *FirebaseAuthFunctionBuilder {
	return fb.onOperation(AuthUserDelete, fn)
}

func (fb *FirebaseAuthFunctionBuilder) onOperation(event string, fn FirebaseAuthFunction) *FirebaseAuthFunctionBuilder {
	fb.Event = event
	fb.Handler = fn
	return fb
}
