package gcp

import (
	"context"
	"google.golang.org/api/option"
)

// DefaultPrefix is the default prefixed used by the [[SecretManager]]
const DefaultPrefix = "gcpsecretmanager"

// SecretManager is the Google Cloud Secret Manager Implementation
// https://cloud.google.com/security/products/secret-manager
type SecretManager struct {
	prefix     string
	clientOpts []option.ClientOption
}

func NewSecretManager(prefix string, opts ...option.ClientOption) *SecretManager {
	return &SecretManager{prefix: prefix, clientOpts: opts}
}

func NewDefaultSecretManager() *SecretManager {
	return NewSecretManager(DefaultPrefix)
}

func (manager *SecretManager) Prefix() string {
	return manager.prefix
}

func (manager *SecretManager) Secret(ctx context.Context, key string) (secret string, err error) {
	c, err := Client(ctx, manager.clientOpts...)
	if err != nil {
		return "", err
	}
	defer func() {
		dErr := c.Close()
		if dErr != nil && err == nil {
			err = dErr
		}
	}()

	return fetchSecret(ctx, c, key)
}
