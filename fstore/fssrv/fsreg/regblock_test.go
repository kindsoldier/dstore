package fsreg

import (
    "fmt"
    //"path/filepath"
    "testing"
    //"math/rand"
    "time"

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

    var fileId      int64   = 1
    var batchId     int64   = 2
    var blockId     int64   = 3
    var blockType   string  = "unk"

    descr0 := dscom.NewBlockDescr()

    ver0 := time.Now().UnixNano()
    ver1 := time.Now().UnixNano()
    ver2 := time.Now().UnixNano()
    require.NotEqual(t, ver0, ver1)
    require.NotEqual(t, ver1, ver2)
    require.NotEqual(t, ver0, ver2)

    descr0.FileId       = fileId
    descr0.BatchId      = batchId
    descr0.BlockId      = blockId
    descr0.BlockType    = blockType
    descr0.BlockVer     = ver0

    descr0.UCounter     = 1
    descr0.BlockSize    = 1024
    descr0.DataSize     = 1111

    descr0.HashAlg       = "a2"
    descr0.HashInit      = "a3"
    descr0.HashSum       = "a4"
    descr0.FilePath      = fmt.Sprintf("a/b/c/qwerty%020d", fileId)

    descr0.FStoreId   = 7
    descr0.BStoreId   = 8
    descr0.SavedLoc   = true
    descr0.SavedRem   = true


    var exists bool

    _ = reg.EraseAllBlockDescrs()
    _ = reg.EraseSpecBlockDescr(descr0.FileId, descr0.BatchId, descr0.BlockId, descr0.BlockType, ver0)
    _ = reg.EraseSpecBlockDescr(descr0.FileId, descr0.BatchId, descr0.BlockId, descr0.BlockType, ver1)
    _ = reg.EraseSpecBlockDescr(descr0.FileId, descr0.BatchId, descr0.BlockId, descr0.BlockType, ver2)

    err = reg.AddNewBlockDescr(descr0)
    require.NoError(t, err)

    err = reg.AddNewBlockDescr(descr0)
    require.Error(t, err)

    descr0.BlockVer     = ver1
    err = reg.AddNewBlockDescr(descr0)
    require.NoError(t, err)

    descr0.BlockVer     = ver2
    err = reg.AddNewBlockDescr(descr0)
    require.NoError(t, err)

    descrs, err := reg.ListAllBlockDescrs()
    require.NoError(t, err)
    require.Equal(t, 3, len(descrs))

    exists, descr1, err := reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr0, descr1)


    err = reg.IncSpecBlockDescrUC(1, fileId, batchId, blockId, blockType, ver2)
    require.NoError(t, err)

    err = reg.DecSpecBlockDescrUC(1, fileId, batchId, blockId, blockType, ver2)
    require.NoError(t, err)

    err = reg.DecSpecBlockDescrUC(1, fileId, batchId, blockId, blockType, ver2)
    require.NoError(t, err)

    exists, descr2, err := reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr2.UCounter, int64(1))

    exists, descr3, err := reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr3.UCounter, int64(1))
    descr3.UCounter = descr0.UCounter
    descr3.BlockVer = descr0.BlockVer
    require.Equal(t, descr0, descr3)

    err = reg.EraseSpecBlockDescr(descr0.FileId, descr0.BatchId, descr0.BlockId, descr0.BlockType, ver0)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.EraseSpecBlockDescr(descr0.FileId, descr0.BatchId, descr0.BlockId, descr0.BlockType, ver1)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, false)

    err = reg.IncSpecBlockDescrUC(1, fileId, batchId, blockId, blockType, ver2)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.EraseSpecBlockDescr(descr0.FileId, descr0.BatchId, descr0.BlockId, descr0.BlockType, ver2)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, false)
}
