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

func TestBlock01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.leveldb")
    defer db.Close()
    require.NoError(t, err)

    reg, err := fsreg.NewReg(db)
    require.NoError(t, err)

    var blockSize int64 = 1024
    var blockId int64 = 1
    var batchId int64 = 2
    var fileId  int64 = 3

    block, err := NewBlock(dataDir, reg, blockId, batchId, fileId, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, block, nil)

    dataSize := blockSize + 1
    buffer := make([]byte, dataSize)
    rand.Read(buffer)
    reader := bytes.NewReader(buffer)

    needSize := int64(200)
    wrSize, err := block.Write(reader, needSize)
    require.NoError(t, err)
    require.Equal(t, wrSize, needSize)

    block, err = OpenBlock(dataDir, reg, blockId, batchId, fileId)
    require.NoError(t, err)
    require.NotEqual(t, block, nil)

    readSize, err := block.Read(io.Discard)
    require.NoError(t, err)
    require.Equal(t, wrSize, readSize)

    err = block.Clean()
    require.NoError(t, err)
}
