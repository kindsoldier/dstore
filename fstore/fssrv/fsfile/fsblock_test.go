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


    var written int64
    var dataSize int64 = 100
    data := make([]byte, dataSize + 200)
    rand.Read(data)

    // Create block in not extsts
    block8, _ := NewBlock(reg, rootDir, fileId, batchId, blockId, blockType, blockSize)
    require.NotEqual(t, block8, nil)

    err = block8.Close()
    require.NoError(t, err)

    // Open block for erasing
    block9, _ := OpenBlock(reg, rootDir, fileId, batchId, blockId, blockType)
    if block9 != nil {
        err = block9.Erase()
        require.NoError(t, err)

        err = block9.Close()
        require.NoError(t, err)
    }

    var need int64 = 100
    // Create new block
    block0, err := NewBlock(reg, rootDir, fileId, batchId, blockId, blockType, blockSize)
    defer block0.Close()
    require.NoError(t, err)
    require.NotEqual(t, block0, nil)

    // Write to new block
    reader0 := bytes.NewReader(data)
    written, err = block0.Write(reader0, need)
    require.NoError(t, err)
    require.Equal(t, need, written)
    // Close block
    err = block0.Close()
    require.NoError(t, err)

    // Reopen block
    block1, err := OpenBlock(reg, rootDir, fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    // Read data to buffer
    writer1 := bytes.NewBuffer(make([]byte, 0))
    written, err = block1.Read(writer1)
    require.NoError(t, err)
    require.Equal(t, need, written)
    require.Equal(t, need, int64(len(writer1.Bytes())))
    require.Equal(t, data[0:need], writer1.Bytes())
    // Close block
    err = block1.Close()
    require.NoError(t, err)

    // Reopen block
    block2, err := OpenBlock(reg, rootDir, fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    // Write to block
    reader1 := bytes.NewReader(data)
    written, err = block2.Write(reader1, need)
    require.NoError(t, err)
    require.Equal(t, dataSize, written)
    // Write to block
    reader2 := bytes.NewReader(data)
    written, err = block2.Write(reader2, need)
    require.NoError(t, err)
    require.Equal(t, dataSize, written)
    // Close block
    err = block2.Close()
    require.NoError(t, err)

    // Re-new block
    _, err = NewBlock(reg, rootDir, fileId, batchId, blockId, blockType, blockSize)
    require.Error(t, err)

    // Reopen block
    block3, err := OpenBlock(reg, rootDir, fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    writer3 := bytes.NewBuffer(make([]byte, 0))
    written, err = block3.Read(writer3)
    require.NoError(t, err)
    require.Equal(t, need * 3, int64(len(writer3.Bytes())))


    err = block3.Erase()
    require.NoError(t, err)

    err = block3.Close()
    require.NoError(t, err)

    // Clean all unised blocks
    for {
        exists, descr, err := reg.GetAnyUnusedBlockDescr()
        require.NoError(t, err)
        if !exists {
            break
        }
        block, err := OpenSpecUnusedBlock(reg, rootDir, descr.FileId, descr.BatchId, descr.BlockId,
                                                        descr.BlockType, descr.BlockVer)
        require.NoError(t, err)
        err = block.Erase()
        require.NoError(t, err)
        err = block.Close()
        require.NoError(t, err)
    }
}
