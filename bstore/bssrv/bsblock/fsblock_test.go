/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsblock

import(
    "bytes"
    "math/rand"
    "testing"

    "github.com/stretchr/testify/require"

    "dstore/dscomm/dskvdb"
    "dstore/bstore/bssrv/bsreg"
)

func TestBlock01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)

    reg, err := bsreg.NewReg(db)
    require.NoError(t, err)

    var fileId      int64 = 1
    var batchId     int64 = 2
    var blockType   int64 = 3
    var blockId     int64 = 4
    var blockSize   int64 = 1024 * 1024 * 16

    block, err := NewBlock(dataDir, fileId, batchId, blockType, blockId, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, block, nil)

    dataSize := blockSize
    buffer := make([]byte, dataSize)
    rand.Read(buffer)
    reader := bytes.NewReader(buffer)

    needSize := blockSize - 1
    wrSize, err := block.Write(reader, needSize)
    require.NoError(t, err)
    require.Equal(t, wrSize, needSize)

    descr := block.Descr()

    block, err = OpenBlock(dataDir, descr)
    require.NoError(t, err)
    require.NotEqual(t, block, nil)

    writer := bytes.NewBuffer(nil)

    readSize, err := block.Read(writer, needSize)
    require.NoError(t, err)
    require.Equal(t, wrSize, readSize)
    require.Equal(t, writer.Bytes(), buffer[0:wrSize])

    err = block.Clean()
    require.NoError(t, err)

    err = reg.DeleteBlock(block.fileId, block.batchId, block.blockType, block.blockId)
    require.NoError(t, err)
}
