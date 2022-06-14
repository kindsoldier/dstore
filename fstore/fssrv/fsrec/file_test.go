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

    err = store.SeedUsers()
    assert.NoError(t, err)

    err = store.SeedBStores()
    assert.NoError(t, err)

    fileName := "qwerty.txt"

    data := make([]byte, 10)
    rand.Read(data)

    reader := bytes.NewReader(data)
    dataSize := int64(len(data))

    userName := "admin"
    err = store.SaveFile(userName, fileName, reader, dataSize)
    assert.NoError(t, err)

    return

    writer := bytes.NewBuffer(make([]byte, 0))

    err = store.LoadFile(userName, fileName, writer)
    assert.NoError(t, err)

    assert.Equal(t, data, writer.Bytes())

    err = store.DeleteFile(userName, fileName)
    assert.NoError(t, err)
}
