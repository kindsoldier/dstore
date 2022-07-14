/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "testing"
    "bytes"
    "math/rand"

    "github.com/stretchr/testify/require"

    "dstore/dscomm/dskvdb"
    "dstore/dscomm/dsalloc"
    "dstore/fstore/fssrv/fsreg"
)


func TestFile01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "storedb")
    defer db.Close()
    require.NoError(t, err)

    reg, err := fsreg.NewReg(db)
    require.NoError(t, err)

    idAlloc, err := dsalloc.OpenAlloc(db, []byte("fileIds"))
    require.NoError(t, err)

    store, err := NewStore(dataDir, reg, idAlloc)
    require.NoError(t, err)

    err = store.SeedUsers()
    require.NoError(t, err)

    var dataSize int64 = 1000 * 1000 * 50
    buffer := make([]byte, dataSize)
    rand.Read(buffer)
    reader := bytes.NewReader(buffer)

    loginDescrs, err := reg.ListUsers()
    require.NoError(t, err)
    require.Equal(t, len(loginDescrs), 2)
    goodLogin := loginDescrs[0].Login
    wrongLogin := "blabla"

    fileName := "/qwerty.txt"
    err = store.SaveFile(goodLogin, fileName, reader, dataSize)
    require.NoError(t, err)

    writer1 := bytes.NewBuffer(nil)
    err = store.LoadFile(goodLogin, fileName, writer1)
    require.NoError(t, err)
    require.Equal(t, int64(len(writer1.Bytes())), dataSize)
    require.Equal(t, writer1.Bytes(), buffer)

    writer2 := bytes.NewBuffer(nil)
    err = store.LoadFile(goodLogin, fileName, writer2)
    require.NoError(t, err)
    require.Equal(t, int64(len(writer2.Bytes())), dataSize)
    require.Equal(t, writer2.Bytes(), buffer)

    writer3 := bytes.NewBuffer(nil)
    err = store.LoadFile(wrongLogin, fileName, writer3)
    require.Error(t, err)
    require.Equal(t, int64(len(writer3.Bytes())), int64(0))

}
