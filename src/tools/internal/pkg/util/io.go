package util

import (
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

//ReadWrapper wraps file or uses stdin
type ReadWrapper struct {
	file   *os.File
	reader io.Reader
}

//WriteWrapper wraps file or uses stdout
type WriteWrapper struct {
	file   *os.File
	writer io.Writer
}

//NewReadWrapper tries open file or uses stdin if file is empty
func NewReadWrapper(file string) (*ReadWrapper, error) {
	res := &ReadWrapper{}
	res.reader = os.Stdin
	if file != "" {
		var err error
		log.Printf("Open file " + file)
		res.file, err = os.Open(file)
		if err != nil {
			return nil, errors.Wrapf(err, "Can't read file %s ", file)
		}
		res.reader = res.file
	}
	return res, nil
}

//Read is main reader method
func (rw *ReadWrapper) Read(p []byte) (int, error) {
	return rw.reader.Read(p)
}

//Close closes file if any
func (rw *ReadWrapper) Close() error {
	if rw.file != nil {
		f := rw.file
		rw.file = nil
		return f.Close()
	}
	return nil
}

//NewWriteWrapper tries open file or uses stdout if file is empty
func NewWriteWrapper(file string) (*WriteWrapper, error) {
	res := &WriteWrapper{}
	res.writer = os.Stdout
	if file != "" {
		var err error
		log.Printf("Open file " + file)
		res.file, err = os.Create(file)
		if err != nil {
			return nil, errors.Wrapf(err, "Can't create file %s ", file)
		}
		res.writer = res.file
	}
	return res, nil
}

//Write is main writer method
func (ww *WriteWrapper) Write(p []byte) (int, error) {
	return ww.writer.Write(p)
}

//WriteString implement StringWritter
func (ww *WriteWrapper) WriteString(s string) (int, error) {
	return io.WriteString(ww.writer, s)
}

//Close closes file if any
func (ww *WriteWrapper) Close() error {
	if ww.file != nil {
		f := ww.file
		ww.file = nil
		return f.Close()
	}
	return nil
}
