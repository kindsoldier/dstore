package fsfile

import (
    "bytes"
    "math/rand"
    "testing"
    "ndstore/fstore/fssrv/fsreg"
    "github.com/stretchr/testify/require"
)

func xxxTest_Block_WriteRead(t *testing.T) {
    var err error

    rootDir := t.TempDir()

    dbPath := "postgres://test@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    var fileId      int64   = 1
    var batchId     int64   = 2
    var blockId     int64   = 3
    var blockType   string  = "tmp"
    var blockSize   int64   = 1024 * 1024

    block0, err := NewBlock(reg, rootDir, fileId, batchId, blockId, blockType, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, block0, nil)

    var written int64
    var dataSize int64 = 100 * 1000
    data := make([]byte, dataSize + 200)
    rand.Read(data)

    reader0 := bytes.NewReader(data)
    written, err = block0.Write(reader0, dataSize)
    require.NoError(t, err)
    require.Equal(t, dataSize, written)

    err = block0.Close()
    require.NoError(t, err)

    block1, err := OpenBlock(reg, rootDir, fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    reader1 := bytes.NewReader(data)
    written, err = block1.Write(reader1, dataSize)
    require.NoError(t, err)
    require.Equal(t, dataSize, written)

    reader2 := bytes.NewReader(data)
    written, err = block1.Write(reader2, dataSize)
    require.NoError(t, err)
    require.Equal(t, dataSize, written)

    err = block1.Close()
    require.NoError(t, err)

    _, err = NewBlock(reg, rootDir, fileId, batchId, blockId, blockType, blockSize)
    require.Error(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    block3, err := OpenBlock(reg, rootDir, fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    written, err = block3.Read(writer)
    require.NoError(t, err)
    require.Equal(t, dataSize * 3, int64(len(writer.Bytes())))

    err = block3.Erase()
    require.NoError(t, err)

    err = block3.Close()
    require.NoError(t, err)
}
