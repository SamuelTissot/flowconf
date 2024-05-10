package flowconf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"reflect"
	"strings"
)

type Builder struct {
	sources  []*StaticSource
	managers []SecretManager
}

func NewBuilder(staticSources ...*StaticSource) *Builder {
	return &Builder{sources: staticSources}
}

func (builder *Builder) Build(config any) error {
	return builder.BuildCtx(context.Background(), config)
}

func (builder *Builder) SetSecretManagers(managers ...SecretManager) {
	builder.managers = managers
}

func (builder *Builder) BuildCtx(ctx context.Context, config any) error {
	err := checkIfConfigIsValid(config)
	if err != nil {
		return err
	}

	err = buildFromSources(config, builder.sources)
	if err != nil {
		return err
	}

	if len(builder.managers) > 0 {
		return resolveSecrets(ctx, config, builder.managers)
	}

	return nil
}

func buildFromSources(config any, sources []*StaticSource) error {
	var err error
	for _, source := range sources {

		switch source.format {
		case Toml:
			err = parseTOML(config, source.reader)
		case Json:
			err = parseJSON(config, source.reader)
		default:
			err = fmt.Errorf("unsupported format: %s", source.format)
		}

		if err != nil {
			return fmt.Errorf("failed to process source: %s, %w", source.name, err)
		}

		err = source.reader.Close()
		if err != nil {
			return fmt.Errorf("failed to close source: %s, %s", source.name, err)
		}
	}

	return nil
}

func parseTOML(config any, r io.Reader) error {
	_, err := toml.NewDecoder(r).Decode(config)
	return err
}

func parseJSON(config any, r io.Reader) error {
	return json.NewDecoder(r).Decode(config)
}

func resolveSecrets(ctx context.Context, config any, managers []SecretManager) error {
	strConfig, err := configToJSON(config)
	if err != nil {
		return err
	}

	subs := findSubstitutions(strConfig)
	for _, sub := range subs {
		manager, err := getManagerForPrefix(sub.managerPrefix, managers)
		if err != nil {
			return err
		}

		secret, err := manager.Secret(ctx, sub.managerKey)
		if err != nil {
			return err
		}

		escapedSecret := escape(secret)
		strConfig = strings.Replace(strConfig, sub.value, escapedSecret, 1)
	}

	return parseJSON(config, strings.NewReader(strConfig))
}

func escape(str string) string {
	str = strings.ReplaceAll(str, "\n", "\\n")
	return strings.ReplaceAll(str, `"`, `\"`)
}

func configToJSON(config any) (string, error) {
	confb, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf(
			"failed to marshal config in order to resolve secrets, %w", err,
		)
	}

	return string(confb), nil
}

type substitutions struct {
	value         string
	managerPrefix string
	managerKey    string
}

func findSubstitutions(str string) []substitutions {
	var out []substitutions

	matches := managerReg.FindAllStringSubmatch(str, -1)
	for i := range matches {
		if len(matches[i]) != 3 {
			continue
		}
		if strings.HasPrefix(matches[i][0], "\\") {
			continue
		}

		out = append(
			out, substitutions{
				value:         matches[i][0],
				managerPrefix: matches[i][1],
				managerKey:    matches[i][2],
			},
		)
	}

	return out
}

func getManagerForPrefix(prefix string, managers []SecretManager) (SecretManager, error) {
	for _, manager := range managers {
		if manager.Prefix() == prefix {
			return manager, nil
		}
	}

	return nil, fmt.Errorf("manager not implemented for prefix: %s", prefix)
}

func checkIfConfigIsValid(config any) error {
	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Pointer {
		return NotAPtrErr
	}

	if rv.IsNil() {
		return IsNilErr
	}

	return nil
}
