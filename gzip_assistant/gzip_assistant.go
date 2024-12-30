package gzipassistant

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func Decrypt(src []byte) ([]byte, error) {
	reader := bytes.NewReader(src)
	if reader == nil {
		return nil, fmt.Errorf("reader is nil")
	}

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	if gzipReader == nil {
		return nil, fmt.Errorf("gzip.NewReader is nil")
	}

	if err := gzipReader.Close(); err != nil {
		return nil, err
	}

	return io.ReadAll(gzipReader)
}

func Encrypt(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(src); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}
	dst := buf.Bytes()
	if dst == nil {
		return nil, fmt.Errorf("compressedResponse is nil")
	}
	return dst, nil
}
