package test

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type SecretManagerMock struct {
	mock.Mock
}

func (manager *SecretManagerMock) Prefix() string {
	return manager.Called().String(0)
}

func (manager *SecretManagerMock) Secret(ctx context.Context, key string) (string, error) {
	args := manager.Called(ctx, key)
	return args.String(0), args.Error(1)
}
