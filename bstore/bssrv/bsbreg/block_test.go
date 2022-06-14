package bsbreg

import (
    "fmt"
    "path/filepath"
    "testing"
    "math/rand"

    "github.com/stretchr/testify/assert"
)

func Test_BlockDescr_InsertSelectDelete(t *testing.T) {
    var err error

    path := filepath.Join(t.TempDir(), "block.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)


    var fileId      int64   = 1
    var batchId     int64   = 2
    var blockId     int64   = 3
    var blockSize   int64   = 1024
    var dataSize    int64   = 1123

    var blockType   string  = "unk"
    var hashAlg     string  = "a2"
    var hashInit    string  = "a3"
    var hashSum     string  = "a4"


    var filePath    string  = fmt.Sprintf("a/b/c/qwerty%020d", fileId)
    var exists bool

    err = reg.AddBlockDescr(fileId, batchId, blockId, blockSize, dataSize, filePath,
                                                      blockType, hashAlg, hashInit, hashSum)
    assert.NoError(t, err)

    err = reg.AddBlockDescr(fileId, batchId, blockId, blockSize, dataSize, filePath,
                                                      blockType, hashAlg, hashInit, hashSum)
    assert.Error(t, err)

    err = reg.AddBlockDescr(fileId, batchId, blockId, blockSize, dataSize, filePath,
                                                      "hihihi", hashAlg, hashInit, hashSum)
    assert.NoError(t, err)


    exists, err = reg.BlockDescrExists(fileId, batchId, blockId, blockType)
    assert.NoError(t, err)
    assert.Equal(t, exists, true)

    exists, err = reg.BlockDescrExists(fileId, batchId, blockId, "hohoho")
    assert.NoError(t, err)
    assert.Equal(t, exists, false)

    exists, err = reg.BlockDescrExists(fileId, batchId, blockId, "hihihi")
    assert.NoError(t, err)
    assert.Equal(t, exists, true)


    exists, err = reg.BlockDescrExists(fileId + 1, batchId, blockId, blockType)
    assert.NoError(t, err)
    assert.Equal(t, exists, false)


    nFileName, nDataSize, err := reg.GetBlockFilePath(fileId, batchId, blockId, blockType)
    assert.NoError(t, err)
    assert.Equal(t, filePath, nFileName)
    assert.Equal(t, dataSize, nDataSize)


    _, _, err = reg.GetBlockFilePath(fileId, batchId, blockId, "hohoho")
    assert.Error(t, err)

    //fileId  += 1
    //batchId += 1
    //blockId += 1
    blockSize += 1024
    filePath = fmt.Sprintf("x/y/z/qwerty%020d", fileId)
    err = reg.UpdateBlockDescr(fileId, batchId, blockId, blockSize, dataSize, filePath)
    assert.NoError(t, err)

    nFileName, nDataSize, err = reg.GetBlockFilePath(fileId, batchId, blockId, blockType)
    assert.NoError(t, err)
    assert.Equal(t, filePath, nFileName)
    assert.Equal(t, dataSize, nDataSize)

    err = reg.DeleteBlockDescr(fileId, batchId, blockId, "hohoho")
    assert.NoError(t, err)

    err = reg.DeleteBlockDescr(fileId, batchId, blockId, "hihihi")
    assert.NoError(t, err)

    err = reg.DeleteBlockDescr(fileId, batchId, blockId, blockType)
    assert.NoError(t, err)
}

func BenchmarkInsert(b *testing.B) {
    var err error
    path := filepath.Join(b.TempDir(), "blocks.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    assert.NoError(b, err)

    err = reg.MigrateDB()
    assert.NoError(b, err)

    b.ResetTimer()

    const numRange int = 1024

    var blockType   string  = "a1"
    var hashAlg     string  = "a2"
    var hashInit    string  = "a3"
    var hashSum     string  = "a4"

    var blockSize   int64   = 1024
    var dataSize    int64   = 123

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            var fileId      int64   = int64(rand.Intn(numRange))
            var batchId     int64   = int64(rand.Intn(numRange))
            var blockId     int64   = int64(rand.Intn(numRange))
            var filePath    string  = fmt.Sprintf("a/b/c/qwerty%020d", fileId)
            err = reg.AddBlockDescr(fileId, batchId, blockId, blockSize, dataSize, filePath,
                                                      blockType, hashAlg, hashInit, hashSum)
            assert.NoError(b, err)
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
    assert.NoError(b, err)

    err = reg.MigrateDB()
    assert.NoError(b, err)

    var hashAlg     string  = "a2"
    var hashInit    string  = "a3"
    var hashSum     string  = "a4"
    var blockType   string  = "unk"

    const numRange int = 1024 * 10
    var i int64
    for i = 0; i < int64(numRange); i++ {
        var fileId      int64   = i
        var batchId     int64   = 1
        var blockId     int64   = 1
        var blockSize   int64   = 1024
        var dataSize    int64   = 123
        var filePath    string  = fmt.Sprintf("a/b/c/qwerty%020d", fileId)
        err = reg.AddBlockDescr(fileId, batchId, blockId, blockSize, dataSize, filePath,
                                                      blockType, hashAlg, hashInit, hashSum)

        assert.NoError(b, err)
    }

    b.ResetTimer()

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            var fileId      int64   = int64(rand.Intn(numRange))
            var batchId     int64   = 1
            var blockId     int64   = 1
            _, _, err = reg.GetBlockFilePath(fileId, batchId, blockId, blockType)
            assert.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}
