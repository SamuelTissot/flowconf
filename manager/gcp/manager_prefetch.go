package gcp

import (
	"context"
	"google.golang.org/api/option"
)

// PrefetchSecretManager will prefetch all secret based on a filter
// see filter documentation: https://cloud.google.com/secret-manager/docs/filtering
//
// Note that the prefetch Manager fetches the LATEST version of each secret return by the filter
type PrefetchSecretManager struct {
	*SecretManager

	cachedSecrets map[string]string
}

// NewPrefetchSecretManager creates a new instance of PrefetchSecretManager.
// It initializes the SecretManager with the default prefix and provided options.
// Then, it fetches filtered secrets based on the provided context, parent, filter, and options.
// Finally, it returns the initialized PrefetchSecretManager with the fetched secrets.
// Example usage:
//
//	manager, err := NewPrefetchSecretManager(ctx, "projects/*/locations/*", "labels.environment:prd", option.WithCredentialsFile("key.json"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
// NOTE ON THE PARENT
// The parent is the resource name of the project associated with the
// [Secrets][google.cloud.secretmanager.v1.Secret], in the format `projects/*` or `projects/*/locations/*`
//
// NOTE ON FILTER
// see: https://cloud.google.com/secret-manager/docs/filtering
//
// ATTENTION
// this implementation fetches the latest versions of each secrets returned by the filter.
// so this means it APPENDS `/versions/latest` to the secret key
// so in order to have a match the keys in the config also needs to have the `/versions/latest`
// EXAMPLE in config:
// @prefix::projects/my-project-id/secrets/password // !!! will not match
// @prefix::projects/my-project-id/secrets/password/versions/latest // !!! WILL MATCH
//
// this enables overriding specific version and still fetching most of the secret with cache
func NewPrefetchSecretManager(
	ctx context.Context,
	parent string,
	filter string,
	opts ...option.ClientOption) (*PrefetchSecretManager, error) {
	manager := NewSecretManager(DefaultPrefix, opts...)

	cachedSecrets, err := fetchFilteredSecrets(ctx, parent, filter, opts)
	if err != nil {
		return nil, err
	}

	return &PrefetchSecretManager{
		SecretManager: manager, cachedSecrets: cachedSecrets,
	}, nil
}

// SetPrefix sets the prefix for the PrefetchSecretManager instance.
// The prefix is used to filter secrets when prefetching.
// Example usage: manager.SetPrefix("my-project/secrets/")
func (manager *PrefetchSecretManager) SetPrefix(prefix string) {
	manager.SecretManager.prefix = prefix
}

// Secret retrieves the secret value for a given key from the PrefetchSecretManager.
// If the secret is already cached, it returns the cached value.
// Otherwise, it delegates the retrieval to the SecretManager and returns the fetched value.
func (manager *PrefetchSecretManager) Secret(ctx context.Context, key string) (secret string, err error) {
	if s, ok := manager.cachedSecrets[key]; ok {
		return s, nil
	}

	return manager.SecretManager.Secret(ctx, key)
}
