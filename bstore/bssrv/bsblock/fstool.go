/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsblock

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "io"
    "math/rand"
    "path/filepath"
    "time"

    "dstore/dscomm/dserr"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func newFilePath() (string) {
    origin := make([]byte, 256)
    rand.Read(origin)
    hasher := sha256.New()
    hasher.Write(origin)
    hashSum := hasher.Sum(nil)
    hashHex := hex.EncodeToString(hashSum)
    fileName := hashHex
    l1 := string(hashHex[0:1])
    l2 := string(hashHex[1:3])
    l3 := string(hashHex[3:5])
    dirPath := filepath.Join(l1, l2, l3)
    filePath := filepath.Join(dirPath, fileName)
    return filePath
}

func copyData(reader io.Reader, writer io.Writer, size int64) (int64, error) {
    var err error
    var bufSize int64 = 1024 * 8
    var total   int64 = 0
    var remains int64 = size
    buffer := make([]byte, bufSize)

    for {
        if remains == 0 {
            return total, dserr.Err(err)
        }
        if remains < bufSize {
            bufSize = remains
        }
        received, err := reader.Read(buffer[0:bufSize])
        if err != nil {
            return total, dserr.Err(err)
        }
        written, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, dserr.Err(err)
        }
        if written != received {
            err = errors.New("write error")
            return total, dserr.Err(err)
        }
        total += int64(written)
        remains -= int64(written)
    }
    return total, dserr.Err(err)
}
