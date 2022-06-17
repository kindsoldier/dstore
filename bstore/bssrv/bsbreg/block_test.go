package bsbreg

import (
    "fmt"
    "path/filepath"
    "testing"
    "math/rand"
    "github.com/stretchr/testify/require"
    "ndstore/dscom"
)

func Test_BlockDescr_InsertSelectDelete(t *testing.T) {
    var err error

    path := filepath.Join(t.TempDir(), "tmp.block.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    var fileId      int64   = 1
    var batchId     int64   = 2
    var blockId     int64   = 3
    var blockType   string  = "unk"

    descr0 := dscom.NewBlockDescr()

    descr0.FileId      = fileId
    descr0.BatchId     = batchId
    descr0.BlockId     = blockId
    descr0.BlockType   = blockType


    descr0.UCounter     = 1
    descr0.BlockSize    = 1024
    descr0.DataSize     = 1123

    descr0.HashAlg       = "a2"
    descr0.HashInit      = "a3"
    descr0.HashSum       = "a4"
    descr0.FilePath      = fmt.Sprintf("a/b/c/qwerty%020d", fileId)

    var exists bool
    var used bool

    err = reg.AddBlockDescr(descr0)
    require.NoError(t, err)

    err = reg.AddBlockDescr(descr0)
    require.Error(t, err)

    descr7 := dscom.NewBlockDescr()
    *descr7 = *descr0
    descr7.BlockType = "hihihi"
    require.NotEqual(t, descr7, descr0)
    err = reg.AddBlockDescr(descr7)
    require.NoError(t, err)

    exists, used, _, _, err = reg.GetBlockParams(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, used, true)
    require.Equal(t, exists, true)

    exists, used, _, _, err = reg.GetBlockParams(fileId, batchId, blockId, "hohoho")
    require.NoError(t, err)
    require.Equal(t, exists, false)
    require.Equal(t, used, false)

    exists, used, _, _, err = reg.GetBlockParams(fileId, batchId, blockId, "hihihi")
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, used, true)

    exists, used, _, _, err = reg.GetBlockParams(fileId + 1, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, false)
    require.Equal(t, used, false)

    exists, used, nFileName, nDataSize, err := reg.GetBlockParams(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, used, true)
    require.Equal(t, descr0.FilePath, nFileName)
    require.Equal(t, descr0.DataSize, nDataSize)

    err = reg.DecBlockDescrUC(fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    err = reg.DecBlockDescrUC(fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    err = reg.DecBlockDescrUC(fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    exists, blockDescr, err := reg.GetUnusedBlockDescr()
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.NotEqual(t, blockDescr, nil)
    require.Equal(t, blockDescr.UCounter, int64(0))

    _, used, _, _, err = reg.GetBlockParams(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, used, false)

    err = reg.IncBlockDescrUC(fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    exists, used, _, _, err = reg.GetBlockParams(fileId, batchId, blockId, "hohoho")
    require.NoError(t, err)
    require.Equal(t, exists, false)
    require.Equal(t, used, false)

    exists, used, nFileName, nDataSize, err = reg.GetBlockParams(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, used, true)
    require.Equal(t, descr0.FilePath, nFileName)
    require.Equal(t, descr0.DataSize, nDataSize)


    err = reg.EraseBlockDescr(fileId, batchId, blockId, "hohoho")
    require.NoError(t, err)

    err = reg.EraseBlockDescr(fileId, batchId, blockId, "hihihi")
    require.NoError(t, err)

    err = reg.EraseBlockDescr(fileId, batchId, blockId, blockType)
    require.NoError(t, err)

    exists, used, nFileName, nDataSize, err = reg.GetBlockParams(fileId, batchId, blockId, blockType)
    require.NoError(t, err)
    require.Equal(t, nFileName, "")
    require.Equal(t, nDataSize, int64(0))
    require.Equal(t, exists, false)
    require.Equal(t, used, false)

}

func BenchmarkInsert(b *testing.B) {
    var err error
    path := filepath.Join(b.TempDir(), "blocks.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    require.NoError(b, err)

    err = reg.MigrateDB()
    require.NoError(b, err)

    b.ResetTimer()

    const numRange int = 1024

    pBench := func(pb *testing.PB) {
        for pb.Next() {

            var fileId      int64   = int64(rand.Intn(numRange))
            var batchId     int64   = int64(rand.Intn(numRange))
            var blockId     int64   = int64(rand.Intn(numRange))
            var blockType   string  = "unk"

            descr0 := dscom.NewBlockDescr()

            descr0.FileId      = fileId
            descr0.BatchId     = batchId
            descr0.BlockId     = blockId
            descr0.BlockType   = blockType

            descr0.UCounter     = 1
            descr0.BlockSize    = 1024
            descr0.DataSize     = 1123
            descr0.HashAlg       = "a2"
            descr0.HashInit      = "a3"
            descr0.HashSum       = "a4"
            descr0.FilePath      = fmt.Sprintf("a/b/c/qwerty%020d", fileId)

            err = reg.AddBlockDescr(descr0)
            require.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}


func BenchmarkSelect(b *testing.B) {
    var err error
    path := filepath.Join(b.TempDir(), "blocks.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    require.NoError(b, err)

    err = reg.MigrateDB()
    require.NoError(b, err)
    numRange := 1024
    var i int64
    for i = 0; i < int64(numRange); i++ {

            var fileId      int64   = int64(rand.Intn(numRange))
            var batchId     int64   = 1
            var blockId     int64   = 1
            var blockType   string  = "unk"

            descr0 := dscom.NewBlockDescr()

            descr0.FileId      = fileId
            descr0.BatchId     = batchId
            descr0.BlockId     = blockId
            descr0.BlockType    = blockType

            descr0.UCounter     = 1
            descr0.BlockSize    = 1024
            descr0.DataSize     = 1123
            descr0.HashAlg       = "a2"
            descr0.HashInit      = "a3"
            descr0.HashSum       = "a4"
            descr0.FilePath      = fmt.Sprintf("a/b/c/qwerty%020d", fileId)

            err = reg.AddBlockDescr(descr0)
        require.NoError(b, err)
    }

    b.ResetTimer()

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            var fileId      int64   = int64(rand.Intn(numRange))
            var batchId     int64   = 1
            var blockId     int64   = 1
            var blockType   string  = "unk"
            _, _, _, _, err = reg.GetBlockParams(fileId, batchId, blockId, blockType)
            require.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}
