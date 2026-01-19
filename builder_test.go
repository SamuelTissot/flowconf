package flowconf_test

import (
	"testing"
	"time"

	"github.com/SamuelTissot/flowconf"
	"github.com/SamuelTissot/flowconf/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuilder_Build_fromStaticSource(t *testing.T) {
	tests := []struct {
		name    string
		sources []*flowconf.StaticSource
		config  any
		want    any
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "one source toml",
			sources: func() []*flowconf.StaticSource {
				sources, err := flowconf.NewSourcesFromEmbeddedFileSystem(
					test.FileSystem,
					"data/config.toml",
				)
				if err != nil {
					t.Fatalf("failed to load sources from embedded filesystem")
				}
				return sources
			}(),
			config: new(test.Config),
			// values are set in the [[test/data/config.toml]] file
			want: &test.Config{
				MeaningOfLife:   42,
				Cats:            []string{"James", "Bond"},
				Pi:              3.14,
				Perfection:      []int{6, 28, 496, 8128},
				BackToTheFuture: time.Date(1985, 10, 21, 1, 22, 0, 0, time.UTC),
				// no manager here so we should get the value
				Secret: "@managerprefix::projects/id/secrets/name-of-secret",
				Tag:    "this should be resolved with the toml tag",
			},
			wantErr: assert.NoError,
		},
		{
			name: "one source json",
			sources: func() []*flowconf.StaticSource {
				sources, err := flowconf.NewSourcesFromEmbeddedFileSystem(
					test.FileSystem,
					"data/config.json",
				)
				if err != nil {
					t.Fatalf("failed to load sources from embedded filesystem")
				}
				return sources
			}(),
			config: new(test.Config),
			// values are set in the [[test/data/config.json]] file
			want: &test.Config{
				MeaningOfLife:   42,
				Cats:            []string{"James", "Bond"},
				Pi:              3.14,
				Perfection:      []int{6, 28, 496, 8128},
				BackToTheFuture: time.Date(1985, 10, 21, 1, 22, 0, 0, time.UTC),
				// no manager here so we should get the value
				Secret: "@managerprefix::projects/id/secrets/name-of-secret",
				Tag:    "this should be resolved with the json tag",
			},
			wantErr: assert.NoError,
		},
		{
			name:    "error if the config is not a pointer",
			sources: nil,
			config:  test.Config{},
			want:    test.Config{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, flowconf.NotAPtrErr)
			},
		},
		{
			name:    "error if the config is not nil",
			sources: nil,
			config: func() *test.Config {
				var v *test.Config
				return v
			}(),
			want: func() *test.Config {
				var v *test.Config
				return v
			}(),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, flowconf.IsNilErr)
			},
		},
		{
			name: "last source will override previous sources",
			sources: func() []*flowconf.StaticSource {
				sources, err := flowconf.NewSourcesFromEmbeddedFileSystem(
					test.FileSystem,
					"data/config.toml",
					"data/config-override-one.json",
					"data/config-override-two.toml",
				)
				if err != nil {
					t.Fatalf("failed to load sources from embedded filesystem")
				}
				return sources
			}(),
			config: new(test.Config),
			// values are set in the [[test/data/*]] file
			want: &test.Config{
				MeaningOfLife:   42,
				Cats:            []string{"Bob", "Morane"},
				Pi:              3.141592653589793,
				Perfection:      []int{6, 28, 496, 8128},
				BackToTheFuture: time.Date(1985, 10, 21, 1, 22, 0, 0, time.UTC),
				// no manager here so we should get the value
				Secret: "secret in local config file",
				Tag:    "this should be resolved with the toml tag",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				builder := flowconf.NewBuilder(tt.sources...)

				err := builder.Build(tt.config)
				tt.wantErr(t, err)
				assert.EqualValues(t, tt.want, tt.config)
			},
		)
	}
}

func TestBuilder_Build_fetchesSecretsFromSecretsManagers(t *testing.T) {
	// /////////////////////// GIVEN ///////////////////////
	var (
		configSourceFile = "data/config.toml"
		prefix           = "managerprefix"                      // this prefix needs to be the same as in the config file [[test/data/config.toml]]
		secretKey        = "projects/id/secrets/name-of-secret" // needs to be the same as in [[test/data/config.toml]]
		secret           = "some very secret secret"
		managerMock      = new(test.SecretManagerMock)
		config           = new(test.Config)
		want             = &test.Config{
			MeaningOfLife:   42,
			Cats:            []string{"James", "Bond"},
			Pi:              3.14,
			Perfection:      []int{6, 28, 496, 8128},
			BackToTheFuture: time.Date(1985, 10, 21, 1, 22, 0, 0, time.UTC),
			// no manager here so we should get the value
			Secret: secret,
			Tag:    "this should be resolved with the toml tag",
		}
	)

	// setup sources
	sources, err := flowconf.NewSourcesFromEmbeddedFileSystem(
		test.FileSystem,
		configSourceFile,
	)
	assert.NoError(t, err)

	// setup mock
	managerMock.On("Prefix").Return(prefix).Once()
	managerMock.On("Secret", mock.Anything, secretKey).Return(secret, nil)

	// setup builder
	builder := flowconf.NewBuilder(sources...)
	builder.SetSecretManagers(managerMock)

	// /////////////////////// WHEN ///////////////////////
	err = builder.Build(config)

	// /////////////////////// THEN ///////////////////////
	assert.NoError(t, err)
	assert.EqualValues(t, want, config)
	managerMock.AssertExpectations(t)
}
