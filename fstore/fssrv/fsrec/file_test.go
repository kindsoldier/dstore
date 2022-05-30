/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec


import (
    "bytes"
    "math/rand"
    "testing"

    "github.com/stretchr/testify/assert"

    "ndstore/fstore/fssrv/fsreg"
)


func Test_File_SaveLoadDelete(t *testing.T) {
    var err error

    baseDir := t.TempDir()

    dbPath := "postgres://pgsql@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    store := NewStore(baseDir, reg)
    fileName := "qwerty.txt"

    data := make([]byte, 1024 * 1024 * 2)
    rand.Read(data)

    reader := bytes.NewReader(data)
    dataSize := int64(len(data))

    err = store.SaveFile(fileName, reader, dataSize)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = store.LoadFile(fileName, writer)
    assert.NoError(t, err)

    assert.Equal(t, data, writer.Bytes())

    err = store.DeleteFile(fileName)
    assert.NoError(t, err)
}
