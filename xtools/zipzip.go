/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package xtools

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io/ioutil"
)

func Zip(origBytes []byte) ([]byte, error) {
	var err error
	var zipBuffer bytes.Buffer
	zipWriter, err := gzip.NewWriterLevel(&zipBuffer, flate.BestCompression)
	if err != nil {
		return zipBuffer.Bytes(), err
	}
	zipWriter.Write(origBytes)
	zipWriter.Flush()
	zipWriter.Close()
	zipBytes := zipBuffer.Bytes()
	return zipBytes, err
}

func Unzip(zipBytes []byte) ([]byte, error) {
	var err error
	origBytes := make([]byte, 0)
	zipBuffer := bytes.NewBuffer(zipBytes)

	zipReader, err := gzip.NewReader(zipBuffer)
	if err != nil {
		return origBytes, err
	}
	origBytes, err = ioutil.ReadAll(zipReader)
	if err != nil {
		return origBytes, err
	}
	return origBytes, err
}
