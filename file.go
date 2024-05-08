package flowconf

import (
	"io"
	"io/fs"
	"os"
)

// Opener is an interface for opening files on the file system or embedded filesystem
type Opener interface {
	Open(filepath string) (io.ReadCloser, error)
}

// osOpener is a wrapper to open OS files
type osOpener func(name string) (*os.File, error)

func (opener osOpener) Open(filepath string) (io.ReadCloser, error) {
	return opener(filepath)
}

// embeddedOpener is a wrapper to open embedded file
type embeddedOpener func(name string) (fs.File, error)

func (opener embeddedOpener) Open(name string) (io.ReadCloser, error) {
	return opener(name)
}
