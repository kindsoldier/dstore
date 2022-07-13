/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import(
    "bytes"
    "math/rand"
    "testing"
    "io"

    "github.com/stretchr/testify/require"

    "dstore/dskvdb"
    "dstore/fstore/fssrv/fsreg"
)

func TestBatch01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)

    reg, err := fsreg.NewReg(db)
    require.NoError(t, err)

    var batchSize int64 = 5
    var blockSize int64 = 1024 * 1024
    var batchId int64 = 2
    var fileId  int64 = 3

    batch, err := NewBatch(reg, dataDir, batchId, fileId, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, batch, nil)

    dataSize := batchSize * blockSize + 1
    buffer := make([]byte, dataSize)
    rand.Read(buffer)
    reader := bytes.NewReader(buffer)

    needSize := int64(batchSize * blockSize - 1)
    wrSize, err := batch.Write(reader, needSize)
    require.NoError(t, err)
    require.Equal(t, needSize, wrSize)

    batch, err = OpenBatch(reg, dataDir, batchId, fileId)
    require.NoError(t, err)
    require.NotEqual(t, batch, nil)

    readSize, err := batch.Read(io.Discard, needSize)
    require.NoError(t, err)
    require.Equal(t, wrSize, readSize)

    err = batch.Clean()
    require.NoError(t, err)
}
