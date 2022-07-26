/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

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
    fileName := hashHex + ".block"
    l1 := string(hashHex[0:1])
    l2 := string(hashHex[1:2])
    dirPath := filepath.Join(l1, l2)
    filePath := filepath.Join(dirPath, fileName)
    return filePath
}

func copyData(reader io.Reader, writer io.Writer, size int64) (int64, bool, error) {
    var err error
    var bufSize int64 = 1024 * 16
    var total   int64 = 0
    var remains int64 = size
    var eof     bool  = false
    buffer := make([]byte, bufSize)

    for {
        if remains == 0 {
            return total, eof,dserr.Err(err)
        }
        if remains < bufSize {
            bufSize = remains
        }
        received, err := reader.Read(buffer[0:bufSize])
        if err == io.EOF {
            eof = true
            err = nil
            return total, eof,dserr.Err(err)
        }
        if err != nil {
            return total, eof,dserr.Err(err)
        }
        written, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, eof,dserr.Err(err)
        }
        if written != received {
            err = errors.New("write error")
            return total, eof,dserr.Err(err)
        }
        total += int64(written)
        remains -= int64(written)
    }
    return total, eof, dserr. Err(err)
}
