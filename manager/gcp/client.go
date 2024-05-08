package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"google.golang.org/api/option"
)

// Client creates a new gcp Secret Manager client based on gRPC
// The returned client must be Closed when it is done being used to clean up its underlying connections.
//
// Note: it's a variable, so it can be extended
var Client = func(ctx context.Context, opts ...option.ClientOption) (*secretmanager.Client, error) {
	c, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to instanciate secret manager client, %w", err)
	}

	return c, err
}
