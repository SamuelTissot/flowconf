package gcp

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func fetchSecret(ctx context.Context, c *secretmanager.Client, key string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: key,
	}
	resp, err := c.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secrest: %s, %w", key, err)
	}

	return string(resp.GetPayload().GetData()), nil
}

func fetchFilteredSecrets(
	ctx context.Context,
	parent string,
	filter string,
	opts []option.ClientOption,
) (cachedSecrets map[string]string, err error) {
	c, err := Client(ctx, opts...)
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

	for {
		s, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get next iterator value, %w", err)
		}

		key := secretLatestVersionPath(s.GetName())

		secret, err := fetchSecret(ctx, c, key)

		cachedSecrets[key] = secret
	}

	return cachedSecrets, nil
}

func secretLatestVersionPath(base string) string {
	return fmt.Sprintf("%s/versions/latest", base)
}
