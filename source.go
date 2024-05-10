package flowconf

import (
	"embed"
	"fmt"
	"io"
	"os"
	"strings"
)

// Format represents a string type that specifies a format.
type Format string

const (
	unknown Format = "unknown"
	Toml    Format = "toml"
	Json    Format = "json"
)

// StaticSource represent a config input
type StaticSource struct {
	name   string
	format Format
	reader io.ReadCloser
}

func NewSource(name string, format Format, reader io.ReadCloser) *StaticSource {
	return &StaticSource{
		name:   name,
		format: format,
		reader: reader,
	}
}

func NewSourcesFromFilepaths(filepaths ...string) ([]*StaticSource, error) {
	return LoadSourcesWithOpener(osOpener(os.Open), filepaths...)
}

func NewSourcesFromEmbeddedFileSystem(fs embed.FS, filepaths ...string) ([]*StaticSource, error) {
	return LoadSourcesWithOpener(embeddedOpener(fs.Open), filepaths...)
}

func LoadSourcesWithOpener(opener Opener, filepaths ...string) ([]*StaticSource, error) {
	var sources []*StaticSource
	for _, filepath := range filepaths {
		f, err := opener.Open(filepath)
		if err != nil {
			return nil, err
		}
		sources = append(sources, NewSource(filepath, detectFormat(filepath), f))

		if err != nil {
			return nil, fmt.Errorf("failed to close opener, %w", err)
		}
	}

	return sources, nil
}

func detectFormat(str string) Format {
	if strings.HasSuffix(str, ".toml") {
		return Toml
	}

	if strings.HasSuffix(str, ".json") {
		return Json
	}

	return unknown
}
