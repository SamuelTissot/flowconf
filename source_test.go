package flowconf

import (
	"bytes"
	"flowconf/test"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNewSourceFromFilepaths_returnTwoSourceForTwoInputFiles(t *testing.T) {
	// /////////////////////// GIVEN ///////////////////////
	var (
		fileOne    = "/some/path/conf.toml"
		fileTwo    = "/other/path/conf.json"
		readerOne  = io.NopCloser(bytes.NewReader([]byte("content one")))
		readerTwo  = io.NopCloser(bytes.NewReader([]byte("content two")))
		openerMock = new(test.FileOpenerMock)
		want       = []*StaticSource{
			NewSource(fileOne, Toml, readerOne),
			NewSource(fileTwo, Json, readerTwo),
		}
	)

	openerMock.On("Open", fileOne).Return(readerOne, nil).Once()
	openerMock.On("Open", fileTwo).Return(readerTwo, nil).Once()

	// /////////////////////// WHEN ///////////////////////
	got, err := LoadSourcesWithOpener(openerMock, fileOne, fileTwo)

	// /////////////////////// THEN ///////////////////////
	assert.NoError(t, err)
	assert.EqualValues(t, want, got)

}
