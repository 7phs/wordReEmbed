package storages

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type TextReader struct {
	reader    io.ReadCloser
	bufReader *bufio.Reader
	delimiter byte
}

func OpenTextReader(fileName string, delimiter byte) (*TextReader, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("failed to open file '" + fileName + "': " + err.Error())
	}

	return &TextReader{
		reader:    file,
		bufReader: bufio.NewReader(file),
		delimiter: delimiter,
	}, nil
}

func (o *TextReader) Close() {
	o.reader.Close()
}

func (o *TextReader) ForEach(cb func(slice []byte) bool) {
	var (
		err   error
		slice []byte
	)

	for err != io.EOF {
		if slice, err = o.ReadSlice(); err == nil || err == io.EOF {
			if !cb(slice) {
				return
			}
		}
	}
}

func (o *TextReader) ReadSlice() ([]byte, error) {
	slice, err := o.bufReader.ReadSlice(o.delimiter)
	if err != nil && err != io.EOF {
		return slice, err
	}

	if err == nil {
		slice = slice[:len(slice)-1]
	}

	return slice, err
}
