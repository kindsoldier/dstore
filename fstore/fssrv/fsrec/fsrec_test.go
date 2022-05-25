/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec


import (
    "bytes"
    "math/rand"
    "testing"
    "io"

    //"ndstore/fstore/fsapi"
    //"ndstore/fstore/fssrv/fsrec"
    //"ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)


func TestSaveLoadDelete(t *testing.T) {
    var err error
    baseDir := "./"
    store := NewStore(baseDir)
    fileName := "qwerty.txt"

    data := make([]byte, 1024 * 16)
    rand.Read(data)

    reader := bytes.NewReader(data)
    dataSize := int64(len(data))

    err = store.SaveFile(fileName, reader, dataSize)
    assert.Equal(t, err, io.EOF)

    err = store.DeleteFile(fileName)
    assert.NoError(t, err)
}
