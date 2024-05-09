package gcp

import (
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"flowconf/test"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"testing"
)

func Test_fetchSecret(t *testing.T) {
	type args struct {
		ctx      context.Context
		accessor SecretVersionAccessor
		key      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "returns secret on valid input",
			args: args{
				ctx: context.Background(),
				accessor: func() SecretVersionAccessor {
					accessorMock := new(SecretVersionAccessorMock)
					accessorMock.On(
						"AccessSecretVersion",
						mock.Anything, // context
						mock.MatchedBy(
							func(req *secretmanagerpb.AccessSecretVersionRequest) bool {
								return req.Name == "secret/key/version/latest"
							},
						),
						mock.Anything,
					).Return(
						&secretmanagerpb.AccessSecretVersionResponse{
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte("the secret value"),
							},
						},
						nil, // the error
					)

					return accessorMock
				}(),
				key: "secret/key/version/latest",
			},
			want:    "the secret value",
			wantErr: assert.NoError,
		},
		{
			name: "returns secret on valid input",
			args: args{
				ctx: context.Background(),
				accessor: func() SecretVersionAccessor {
					accessorMock := new(SecretVersionAccessorMock)
					accessorMock.On(
						"AccessSecretVersion",
						mock.Anything, // context
						mock.MatchedBy(
							func(req *secretmanagerpb.AccessSecretVersionRequest) bool {
								return req.Name == "secret/key/version/latest"
							},
						),
						mock.Anything,
					).Return(
						&secretmanagerpb.AccessSecretVersionResponse{
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte("the secret value"),
							},
						},
						nil, // the error
					)

					return accessorMock
				}(),
				key: "secret/key/version/latest",
			},
			want:    "the secret value",
			wantErr: assert.NoError,
		},
		{
			name: "returns an error when the secret does not exist",
			args: args{
				ctx: context.Background(),
				accessor: func() SecretVersionAccessor {
					accessorMock := new(SecretVersionAccessorMock)
					accessorMock.On(
						"AccessSecretVersion",
						mock.Anything, // context
						mock.MatchedBy(
							func(req *secretmanagerpb.AccessSecretVersionRequest) bool {
								return req.Name == "secret/key/version/latest"
							},
						),
						mock.Anything,
					).Return(
						nil,
						test.ExpectedErr, // the error
					)

					return accessorMock
				}(),
				key: "secret/key/version/latest",
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorIs(t, err, test.ExpectedErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := fetchSecret(tt.args.ctx, tt.args.accessor, tt.args.key)
				tt.wantErr(t, err)
				assert.Equalf(t, tt.want, got, "fetchSecret did not return expected results")
			},
		)
	}
}

func Test_fetchFilteredSecrets_happyPath(t *testing.T) {
	// /////////////////////// GIVEN ///////////////////////
	var (
		ctx                = context.Background()
		parent             = "project/parent"
		filter             = "labels.environment:prd"
		iteratorMock       = new(IteratorMock)
		clientMock         = new(ClientMock)
		opts               = []gax.CallOption(nil)
		secretOne          = &secretmanagerpb.Secret{Name: "secret/key/one"}
		secretTwo          = &secretmanagerpb.Secret{Name: "secret/key/two"}
		accessSecretOneReq = &secretmanagerpb.AccessSecretVersionRequest{
			Name: "secret/key/one/versions/latest",
		}
		accessSecretOneResponse = &secretmanagerpb.AccessSecretVersionResponse{
			Name: "secret/key/one/versions/latest",
			Payload: &secretmanagerpb.SecretPayload{
				Data: []byte("the secret ONE value"),
			},
		}
		accessSecretTwoReq = &secretmanagerpb.AccessSecretVersionRequest{
			Name: "secret/key/two/versions/latest",
		}
		accessSecretTwoResponse = &secretmanagerpb.AccessSecretVersionResponse{
			Name: "secret/key/two/versions/latest",
			Payload: &secretmanagerpb.SecretPayload{
				Data: []byte("the secret TWO value"),
			},
		}

		listSecretsReq = &secretmanagerpb.ListSecretsRequest{
			Parent: parent,
			Filter: filter,
		}
		want = map[string]string{
			"secret/key/one/versions/latest": "the secret ONE value",
			"secret/key/two/versions/latest": "the secret TWO value",
		}
	)

	// setup iterator
	iteratorMock.On("Next").Return(secretOne, nil).Once()
	iteratorMock.On("Next").Return(secretTwo, nil).Once()
	iteratorMock.On("Next").Return(nil, iterator.Done)
	// setup client
	clientMock.On("ListSecrets", ctx, listSecretsReq, opts).Return(iteratorMock)
	clientMock.On("Close").Return(nil).Once()
	clientMock.On("AccessSecretVersion", ctx, accessSecretOneReq, opts).Return(accessSecretOneResponse, nil).Once()
	clientMock.On("AccessSecretVersion", ctx, accessSecretTwoReq, opts).Return(accessSecretTwoResponse, nil).Once()

	// mockeyPatch the NewClient Function to return the clientMock
	newClientStub := func(ctx context.Context, opts ...option.ClientOption) (Client, error) {
		return clientMock, nil
	}
	defer test.MonkeyPatch(&NewClient, &newClientStub)()

	// /////////////////////// WHEN ///////////////////////
	got, err := fetchFilteredSecrets(ctx, parent, filter, nil)

	// /////////////////////// THEN ///////////////////////
	assert.NoError(t, err)
	assert.EqualValues(t, want, got)
}

// **************************************************************************
// * MOCKS
// **************************************************************************

type SecretVersionAccessorMock struct {
	mock.Mock
}

func (s *SecretVersionAccessorMock) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	args := s.Called(ctx, req, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*secretmanagerpb.AccessSecretVersionResponse), args.Error(1)
}

type ClientMock struct {
	mock.Mock
}

func (client *ClientMock) Close() error {
	return client.Called().Error(0)
}

func (client *ClientMock) AccessSecretVersion(
	ctx context.Context,
	req *secretmanagerpb.AccessSecretVersionRequest,
	opts ...gax.CallOption,
) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	args := client.Called(ctx, req, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*secretmanagerpb.AccessSecretVersionResponse), args.Error(1)
}

func (client *ClientMock) ListSecrets(
	ctx context.Context,
	req *secretmanagerpb.ListSecretsRequest,
	opts ...gax.CallOption,
) SecretIterator {
	return client.Called(ctx, req, opts).Get(0).(SecretIterator)
}

type IteratorMock struct {
	mock.Mock
}

func (i *IteratorMock) Next() (*secretmanagerpb.Secret, error) {
	args := i.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*secretmanagerpb.Secret), args.Error(1)
}
