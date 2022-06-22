package fsfile

import (
    "bytes"
    "math/rand"
    "testing"
    "ndstore/fstore/fssrv/fsreg"
    "github.com/stretchr/testify/require"
)

func xxxTest_Batch_WriteRead_Size(t *testing.T) {
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
    var batchSize   int64   = 5
    var blockSize   int64   = 1024 * 1024

    batch0, err := NewBatch(reg, rootDir, fileId, batchId, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, batch0, nil)

    var written int64
    var dataSize int64 = batchSize * blockSize
    data := make([]byte, dataSize + 200)
    rand.Read(data)

    reader0 := bytes.NewReader(data)
    written, err = batch0.Write(reader0, blockSize + 1)
    require.NoError(t, err)
    require.Equal(t, int64(blockSize + 1), written)

    err = batch0.Close()
    require.NoError(t, err)

    batch1, err := OpenBatch(reg, rootDir, fileId, batchId)
    require.NoError(t, err)

    reader1 := bytes.NewReader(data)
    written, err = batch1.Write(reader1, 100)
    require.NoError(t, err)
    require.Equal(t, int64(100), written)

    reader2 := bytes.NewReader(data)
    written, err = batch1.Write(reader2, blockSize + 4)
    require.NoError(t, err)
    require.Equal(t, blockSize + 4, written)

    err = batch1.Close()
    require.NoError(t, err)


    _, err = NewBatch(reg, rootDir, fileId, batchId, batchSize, blockSize)
    require.Error(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    batch3, err := OpenBatch(reg, rootDir, fileId, batchId)
    require.NoError(t, err)

    written, err = batch3.Read(writer)
    require.NoError(t, err)
    require.Equal(t, int64(blockSize + 1 + 100 + blockSize + 4 ), int64(len(writer.Bytes())))

    err = batch3.Erase()
    require.NoError(t, err)

    err = batch3.Close()
    require.NoError(t, err)
}

func xxxTest_Batch_WriteRead_Data(t *testing.T) {
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
    var batchSize   int64   = 5
    var blockSize   int64   = 1024 * 1024

    batch0, err := NewBatch(reg, rootDir, fileId, batchId, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, batch0, nil)

    var written int64
    var dataSize int64 = batchSize * blockSize
    data := make([]byte, dataSize + 200)
    rand.Read(data)

    reader0 := bytes.NewReader(data)
    need := blockSize * 3 + 1
    written, err = batch0.Write(reader0, need)
    require.NoError(t, err)
    require.Equal(t, need, written)

    err = batch0.Close()
    require.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    batch3, err := OpenBatch(reg, rootDir, fileId, batchId)
    require.NoError(t, err)

    read, err := batch3.Read(writer)
    require.NoError(t, err)
    require.Equal(t, need, read)
    require.Equal(t, need, int64(len(writer.Bytes())))
    require.Equal(t, data[0:need], writer.Bytes())

    err = batch3.Erase()
    require.NoError(t, err)

    err = batch3.Close()
    require.NoError(t, err)
}
