package fsfile

import (
    "bytes"
    "math/rand"
    "testing"
    "ndstore/fstore/fssrv/fsreg"
    "github.com/stretchr/testify/require"
)

func xxxTest_Batch_WriteRead(t *testing.T) {
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
    var blockSize   int64   = 1024

    // Open batch and erase
    batch8, err := OpenBatch(reg, rootDir, fileId, batchId)
    if err == nil && batch8 != nil {
        err = batch8.Erase()
        require.NoError(t, err)

        err = batch8.Close()
        require.NoError(t, err)
    }
    // Prepare data
    var dataSize int64 = batchSize * blockSize + 1
    data := make([]byte, dataSize)
    rand.Read(data)

    // Create batch
    batch0, err := NewBatch(reg, rootDir, fileId, batchId, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, batch0, nil)
    // Write to batch
    var need int64 = blockSize + 1
    reader0 := bytes.NewReader(data)
    written0, err := batch0.Write(reader0, need)
    require.NoError(t, err)
    require.Equal(t, need, written0)
    // Close batch
    err = batch0.Close()
    require.NoError(t, err)

    // New batch with the same parameters
    _, err = NewBatch(reg, rootDir, fileId, batchId, batchSize, blockSize)
    require.Error(t, err)

    // Reopen batch
    batch1, err := OpenBatch(reg, rootDir, fileId, batchId)
    require.NoError(t, err)
    // Read data
    writer1 := bytes.NewBuffer(make([]byte, 0))
    written1, err := batch1.Read(writer1)
    require.NoError(t, err)
    require.Equal(t, need, written1)
    require.Equal(t, need, int64(len(writer1.Bytes())))
    require.Equal(t, data[0:need], writer1.Bytes())
    // Write yet data
    reader1 := bytes.NewReader(data)
    written1, err = batch1.Write(reader1, need)
    require.NoError(t, err)
    require.Equal(t, need, written1)
    // Write yet data
    reader2 := bytes.NewReader(data)
    written1, err = batch1.Write(reader2, need)
    require.NoError(t, err)
    require.Equal(t, need, written1)
    // Close batch
    err = batch1.Close()
    require.NoError(t, err)

    // Open batch
    batch3, err := OpenBatch(reg, rootDir, fileId, batchId)
    require.NoError(t, err)
    // Read from batch
    writer3 := bytes.NewBuffer(make([]byte, 0))
    read3, err := batch3.Read(writer3)
    require.NoError(t, err)
    require.Equal(t, need * 3, read3)
    require.Equal(t, need * 3, int64(len(writer3.Bytes())))
    // Check data
    require.Equal(t, data[0:need], writer3.Bytes()[0:need])
    require.Equal(t, data[0:need], writer3.Bytes()[0+need:need+need])
    require.Equal(t, data[0:need], writer3.Bytes()[0+need*2:need+need*2])
    // Erase batch
    err = batch3.Erase()
    require.NoError(t, err)
    // Close batch
    err = batch3.Close()
    require.NoError(t, err)

    // Open erased batch
    _, err = OpenBatch(reg, rootDir, fileId, batchId)
    require.Error(t, err)
}
