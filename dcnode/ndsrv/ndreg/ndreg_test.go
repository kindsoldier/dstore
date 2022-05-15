package ndreg

import (
    "fmt"
    "path/filepath"
    "testing"
    "math/rand"

    "github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
    var err error
    path := filepath.Join(t.TempDir(), "reg.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    var count int64 = 128
    var i int64
    for i = 0; i < count; i++ {
        var clusterId   int64   = 1
        var fileId      int64   = 2
        var batchId     int64   = 3
        var blockId     int64   = i
        var blockSize   int64   = 1024
        var filePath    string  = fmt.Sprintf("a/b/c/qwerty%020d", i)

        err = reg.AddBlock(clusterId, fileId, batchId, blockId, blockSize, filePath)
        assert.NoError(t, err)
        var exists bool

        exists, err = reg.BlockExists(clusterId, fileId, batchId, blockId)
        assert.NoError(t, err)
        assert.Equal(t, exists, true)

        exists, err = reg.BlockExists(clusterId, fileId + 1, batchId, blockId)
        assert.NoError(t, err)
        assert.Equal(t, exists, false)

        nFileName, _, err := reg.GetBlock(clusterId, fileId, batchId, blockId)
        assert.NoError(t, err)
        assert.Equal(t, filePath, nFileName)

        //err := reg.DeleteBlock(clusterId, fileId, batchId, blockId)
        //assert.NoError(t, err)
    }
}

func BenchmarkInsert(b *testing.B) {
    var err error
    path := filepath.Join(b.TempDir(), "reg.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    assert.NoError(b, err)

    err = reg.MigrateDB()
    assert.NoError(b, err)

    const numRange int = 1024
    pBench := func(pb *testing.PB) {
        for pb.Next() {
            var clusterId   int64   = int64(rand.Intn(numRange))
            var fileId      int64   = int64(rand.Intn(numRange))
            var batchId     int64   = int64(rand.Intn(numRange))
            var blockId     int64   = int64(rand.Intn(numRange))
            var blockSize   int64   = 1024
            var filePath    string  = fmt.Sprintf("a/b/c/qwerty%020d", fileId)

            err = reg.AddBlock(clusterId, fileId, batchId, blockId, blockSize, filePath)
            assert.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}


func BenchmarkSelect(b *testing.B) {
    var err error
    path := filepath.Join(b.TempDir(), "reg.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    assert.NoError(b, err)

    err = reg.MigrateDB()
    assert.NoError(b, err)


    const numRange int = 1024 * 10
    var i int64
    for i = 0; i < int64(numRange); i++ {
        var clusterId   int64   = 1
        var fileId      int64   = i
        var batchId     int64   = 1
        var blockId     int64   = 1
        var blockSize   int64   = 1024
        var filePath    string  = fmt.Sprintf("a/b/c/qwerty%020d", fileId)
        err = reg.AddBlock(clusterId, fileId, batchId, blockId, blockSize, filePath)
        assert.NoError(b, err)
    }

    b.ResetTimer()

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            var clusterId   int64   = 1
            var fileId      int64   = int64(rand.Intn(numRange))
            var batchId     int64   = 1
            var blockId     int64   = 1
            _, _, err = reg.GetBlock(clusterId, fileId, batchId, blockId)
            assert.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}
