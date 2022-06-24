/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsrpc

import (
    "context"
    "errors"
    "io"
    "fmt"
    "sync"
    "ndstore/dslog"
)

func ReadBytes(reader io.Reader, size int64) ([]byte, error) {
    buffer := make([]byte, size)
    read, err := io.ReadFull(reader, buffer)
    return buffer[0:read], err
}

func CopyBytes(reader io.Reader, writer io.Writer, dataSize int64) (int64, error) {
    var err error
    var bSize int64 = 1024 * 4
    var total int64 = 0
    var remains int64 = dataSize
    buffer := make([]byte, bSize)

    for {
        if reader == nil {
            return total, errors.New("reader is nil")
        }
        if writer == nil {
            return total, errors.New("writer is nil")
        }
        if remains == 0 {
            return total, err
        }
        if remains < bSize {
            bSize = remains
        }
        received, err := reader.Read(buffer[0:bSize])
        if err != nil {
            return total, fmt.Errorf("read error: %v", err)
        }
        recorded, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, fmt.Errorf("write error: %v", err)
        }
        if recorded != received {
            return total, errors.New("size mismatch")
        }
        total += int64(recorded)
        remains -= int64(recorded)
    }
    return total, err
}
