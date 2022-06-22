package fsreg

import (
    //"fmt"
    //"path/filepath"
    "testing"
    //"math/rand"

    "github.com/stretchr/testify/require"
    "ndstore/dscom"
)

func Test_BlockDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    descr0 := dscom.NewBlockDescr()

    var fileId      int64   = 1
    var batchId     int64   = 2
    var blockId     int64   = 3
    var blockType   string  = "tmp"


    descr0.FileId     = fileId
    descr0.BatchId    = batchId
    descr0.BlockId    = blockId
    descr0.BlockType  = blockType

    descr0.BlockSize  = 1024
    descr0.DataSize   = 126
    descr0.FilePath   = "a/b/c/qwerty"
    descr0.HashAlg    = "hway"
    descr0.HashInit   = "hashinit"
    descr0.HashSum    = "hashsum"

    descr0.FStoreId   = 7
    descr0.BStoreId   = 8
    descr0.SavedLoc   = true
    descr0.SavedRem   = true

    var exists bool

    err = reg.EraseBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    err = reg.AddBlockDescr(descr0)
    require.NoError(t, err)

    err = reg.AddBlockDescr(descr0)
    require.Error(t, err)

    err = reg.EraseBlockDescr(fileId + 1, batchId, blockId, blockType)
    require.NoError(t, err)

    descr7 := dscom.NewBlockDescr()
    *descr7 = *descr0
    descr7.FileId += 1
    err = reg.AddBlockDescr(descr7)
    require.NoError(t, err)

    exists, descr1, err := reg.GetBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr0, descr1)

    err = reg.EraseBlockDescr(fileId + 1, batchId, blockId, blockType)
    require.NoError(t, err)

    exists, _, err = reg.GetBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.EraseBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    exists, _, err = reg.GetBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, false)

}
