package flowconf

import (
	"context"
	"regexp"
)

// managerReg is used to find pattern in the configuration to substitute with a secret
var managerReg = regexp.MustCompile(`@(\w+)::([\w/-]+)`)

// SecretManager interface defines the methods required for accessing secrets.
type SecretManager interface {
	// Prefix return the string to look for in the configuration in order to
	// do the substitution in the format of @<PREFIX>::<PATH/KEY>
	//
	// Example for GCP Secret Manager
	// Prefix() --> gcpsecretmanager
	// Usage --> @gcpsecretmanager::project/*/location/*/...
	Prefix() string
	// Secret returns the secret for the given key
	// the key is the @<PREFIX>::<KEY>
	Secret(ctx context.Context, key string) (string, error)
}
