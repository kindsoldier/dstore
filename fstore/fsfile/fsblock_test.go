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

    "dstore/dskeydb"
    "dstore/fsreg"
    "dstore/dslog"
)

func xxTestBlock00(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    var blockSize int64 = 1024
    block, err := NewBlock(dataDir, blockSize)

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

    db, err := dskeydb.OpenKV(dataDir, "tmp.leveldb")
    defer db.Close()
    require.NoError(t, err)

    //alloc, err := fsreg.NewBlockAlloc(db)
    //require.NoError(t, err)

    reg, err := fsreg.NewBlockReg(db)
    require.NoError(t, err)

    var blockSize int64 = 1024
    block, err := NewBlock(dataDir, blockSize)

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

    err = reg.AddBlock(1, descr)
    require.NoError(t, err)

    has, descr, err := reg.GetBlock(1)
    require.NoError(t, err)
    require.Equal(t, has, true)

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
