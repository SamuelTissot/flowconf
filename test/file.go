package test

import (
	"github.com/stretchr/testify/mock"
	"io"
)

type FileOpenerMock struct {
	mock.Mock
}

func (f *FileOpenerMock) Open(file string) (io.ReadCloser, error) {
	args := f.Called(file)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
