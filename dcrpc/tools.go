package dcrpc

import (
    "io"
)

func ReadBytes(reader io.Reader, size int64) ([]byte, error) {
    buffer := make([]byte, size)
    read, err := io.ReadFull(reader, buffer)
    return buffer[0:read], err
}
