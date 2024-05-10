package fake

import (
	"context"
	"fmt"
)

type Manager struct {
	secrets map[string]string
}

func NewManager(secrets map[string]string) *Manager {
	return &Manager{secrets: secrets}
}

func (manager *Manager) Prefix() string {
	return "prefix"
}

func (manager *Manager) Secret(_ context.Context, key string) (string, error) {
	if v, ok := manager.secrets[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("key: %s not found", key)
}
