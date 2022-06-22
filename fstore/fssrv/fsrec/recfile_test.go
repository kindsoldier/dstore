/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec


import (
    "bytes"
    "math/rand"
    "testing"

    "github.com/stretchr/testify/require"

    "ndstore/fstore/fssrv/fsreg"
)


func Test_File_SaveLoadDelete(t *testing.T) {
    var err error

    baseDir := t.TempDir()

    dbPath := "postgres://pgsql@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    store := NewStore(baseDir, reg)

    err = store.SeedUsers()
    require.NoError(t, err)

    err = store.SeedBStores()
    require.NoError(t, err)

    fileName := "qwerty.txt"

    data := make([]byte, 10)
    rand.Read(data)

    reader := bytes.NewReader(data)
    dataSize := int64(len(data))

    userName := "admin"

    //err = store.DeleteFile(userName, fileName)
    //require.NoError(t, err)

    err = store.SaveFile(userName, fileName, reader, dataSize)
    require.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = store.LoadFile(userName, fileName, writer)
    require.NoError(t, err)

    require.Equal(t, data, writer.Bytes())

    return

    err = store.DeleteFile(userName, fileName)
    require.NoError(t, err)
}
