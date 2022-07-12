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
    "dstore/dslog"
)

func TestBlock00(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    var blockSize int64 = 1024
    var blockId int64 = 1
    var batchId int64 = 2
    var fileId  int64 = 3

    block, err := NewBlock(dataDir, blockId, batchId, fileId, blockSize)

    dataSize := blockSize + 1
    buffer := make([]byte, dataSize)
    rand.Read(buffer)
    reader := bytes.NewReader(buffer)

    needSize := int64(200)
    wrSize, err := block.Write(reader, needSize)
    require.NoError(t, err)
    require.Equal(t, wrSize, needSize)

    descr := block.Descr()

    block, err = OpenBlock(dataDir, descr)
    require.NoError(t, err)

    readSize, err := block.Read(io.Discard)
    require.NoError(t, err)
    require.Equal(t, wrSize, readSize)

    err = block.Clean()
    require.NoError(t, err)
}


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


    block, err := NewBlock(dataDir, blockId, batchId, fileId, blockSize)

    descr := block.Descr()
    descrBin, err := descr.Pack()
    require.NoError(t, err)

    dslog.LogDebug(string(descrBin))

    dataSize := blockSize + 1
    buffer := make([]byte, dataSize)
    rand.Read(buffer)
    reader := bytes.NewReader(buffer)

    needSize := int64(200)
    wrSize, err := block.Write(reader, needSize)
    require.NoError(t, err)
    require.Equal(t, wrSize, needSize)

    descr = block.Descr()
    descrBin, err = descr.Pack()
    require.NoError(t, err)

    dslog.LogDebug(string(descrBin))

    err = reg.PutBlock(descr)
    require.NoError(t, err)

    has, err := reg.HasBlock(blockId, batchId, fileId)
    require.NoError(t, err)
    require.Equal(t, has, true)

    descr, err = reg.GetBlock(blockId, batchId, fileId)
    require.NoError(t, err)

    block, err = OpenBlock(dataDir, descr)
    require.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    readSize, err := block.Read(writer)
    require.NoError(t, err)
    require.Equal(t, wrSize, readSize)

    err = block.Clean()
    require.NoError(t, err)

    descr = block.Descr()
    descrBin, err = descr.Pack()
    require.NoError(t, err)

    dslog.LogDebug(string(descrBin))

}
