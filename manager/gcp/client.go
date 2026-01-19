package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"errors"
	"fmt"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
)

// SecretVersionAccessor is an interface for accessing secret versions.
type SecretVersionAccessor interface {
	AccessSecretVersion(
		ctx context.Context,
		req *secretmanagerpb.AccessSecretVersionRequest,
		opts ...gax.CallOption,
	) (*secretmanagerpb.AccessSecretVersionResponse, error)
}

// SecretIterator is an interface for iterating over secrets.
// It defines a single method Next() that returns the next secret and an error.
type SecretIterator interface {
	Next() (*secretmanagerpb.Secret, error)
}

// SecretLister is an interface for listing secrets.
// It defines a single method ListSecrets() that returns a SecretIterator.
// The SecretIterator provides a way to iterate over secrets by calling Next().
// The method takes a context, a ListSecretsRequest, and optional gax.CallOptions.
// It returns a SecretIterator that can be used to retrieve the next secret and an error.
type SecretLister interface {
	ListSecrets(
		ctx context.Context,
		req *secretmanagerpb.ListSecretsRequest,
		opts ...gax.CallOption,
	) SecretIterator
}

// Client is an interface that represents a client for accessing secret versions and listing secrets.
// It extends the io.Closer interface for closing the client connections.
// It also extends the SecretVersionAccessor and SecretLister interfaces.
// The SecretVersionAccessor interface provides a method for accessing secret versions.
// The SecretLister interface provides a method for listing secrets.
type Client interface {
	io.Closer
	SecretVersionAccessor
	SecretLister
}

// ClientWrapper is a type that wraps the secretmanager.Client type.
// It provides a way to extend the functionality of the secretmanager.Client type.
type ClientWrapper struct {
	*secretmanager.Client
}

func NewClientWrapper(c *secretmanager.Client) *ClientWrapper {
	return &ClientWrapper{c}
}

// ListSecrets overloads the ListSecrets method on the client to return the secretIterator interface
func (wrapper *ClientWrapper) ListSecrets(
	ctx context.Context,
	req *secretmanagerpb.ListSecretsRequest,
	opts ...gax.CallOption,
) SecretIterator {
	return wrapper.Client.ListSecrets(ctx, req, opts...)
}

// NewClient creates a new gcp Secret Manager client based on gRPC
// The returned client must be Closed when it is done being used to clean up its underlying connections.
//
// Note: it's a variable, so it can be extended
var NewClient = func(ctx context.Context, opts ...option.ClientOption) (Client, error) {
	c, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to instanciate secret manager client, %w", err)
	}

	return NewClientWrapper(c), err
}

func fetchSecret(ctx context.Context, accessor SecretVersionAccessor, key string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: key,
	}
	resp, err := accessor.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secrets: %s, %w", key, err)
	}

	return string(resp.GetPayload().GetData()), nil
}

func fetchFilteredSecrets(
	ctx context.Context,
	parent string,
	filter string,
	nWorkers int,
	opts []option.ClientOption,
) (cachedSecrets map[string]string, err error) {
	semaphore := make(chan struct{}, nWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errGather := []error{}

	c, err := NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	defer func() {
		dErr := c.Close()
		if dErr != nil && err == nil {
			err = dErr
		}
	}()

	it := c.ListSecrets(
		ctx, &secretmanagerpb.ListSecretsRequest{
			Parent: parent,
			Filter: filter,
		},
	)

	cachedSecrets = map[string]string{}
	for {
		s, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get next iterator value, %w", err)
		}

		wg.Add(1)
		go func(secretName string) {
			defer wg.Done()
			semaphore <- struct{}{}        // make sure we don't flood system
			defer func() { <-semaphore }() // release worker

			key := secretLatestVersionPath(secretName)

			fmt.Printf("fetching secret : %s\n", key)

			secret, err := fetchSecret(ctx, c, key)

			mu.Lock()
			if err != nil {
				errGather = append(errGather, err)
			} else {
				cachedSecrets[key] = secret
			}

			mu.Unlock()
		}(s.GetName())

		wg.Wait()
	}

	return cachedSecrets, nil
}

func secretLatestVersionPath(base string) string {
	return fmt.Sprintf("%s/versions/latest", base)
}
