/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bstore

import (
    "testing"
    "bytes"
    "math/rand"

    "github.com/stretchr/testify/require"

    "dstore/dscomm/dskvdb"
    "dstore/bstore/bssrv/bsreg"
)


func TestBlock01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "storedb")
    defer db.Close()
    require.NoError(t, err)

    reg, err := bsreg.NewReg(db)
    require.NoError(t, err)

    store, err := NewStore(dataDir, reg)
    require.NoError(t, err)

    var dataSize int64 = 1000 * 1000
    buffer := make([]byte, dataSize)
    rand.Read(buffer)
    reader := bytes.NewReader(buffer)

    var fileId      int64 = 1
    var batchId     int64 = 2
    var blockType   int64 = 3
    var blockId     int64 = 4
    var blockSize   int64 = 1024 * 1024 * 16

    err = store.SaveBlock(fileId, batchId, blockType, blockId, blockSize, reader, dataSize)
    require.NoError(t, err)

    writer1 := bytes.NewBuffer(nil)
    err = store.LoadBlock(fileId, batchId, blockType, blockId, writer1, dataSize)
    require.NoError(t, err)
    require.Equal(t, int64(len(writer1.Bytes())), dataSize)
    require.Equal(t, writer1.Bytes(), buffer)

    writer2 := bytes.NewBuffer(nil)
    err = store.LoadBlock(fileId, batchId, blockType, blockId, writer2, dataSize)
    require.NoError(t, err)
    require.Equal(t, int64(len(writer2.Bytes())), dataSize)
    require.Equal(t, writer2.Bytes(), buffer)

    err = store.DeleteBlock(fileId, batchId, blockType, blockId)
    require.NoError(t, err)

    writer3 := bytes.NewBuffer(nil)
    err = store.LoadBlock(fileId, batchId, blockType, blockId, writer3, dataSize)
    require.Error(t, err)
    require.Equal(t, int64(len(writer3.Bytes())), int64(0))

}
