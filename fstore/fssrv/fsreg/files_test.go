package fsreg

import (
    "encoding/json"
    "bytes"
    "io"
    "testing"
    "math/rand"

    "github.com/stretchr/testify/assert"

    "ndstore/fstore/fssrv/fsfile"
)

func Test_FileDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dbPath := "postgres://pgsql@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)


    var batchSize   int64 = 5
    var blockSize   int64 = 1024 * 1024

    baseDir := t.TempDir()

    fileId, err := reg.GetNewFileId()
    assert.NoError(t, err)

    file := fsfile.NewFile(baseDir, fileId, batchSize, blockSize)

    err = file.Open()
    assert.NoError(t, err)

    err = file.Truncate()
    assert.NoError(t, err)

    var dataSize int64 = 6 * (batchSize * blockSize) + blockSize + 2
    data := make([]byte, dataSize)
    rand.Read(data)
    reader := bytes.NewReader(data)

    written, err := file.Write(reader)
    assert.Error(t, err)
    assert.Equal(t, err, io.EOF)
    assert.Equal(t, dataSize, written)

    origFileDescr, err := file.Meta()
    assert.NoError(t, err)

    err = file.Close()
    assert.NoError(t, err)

    //err = reg.DeleteFileDescr(fileId)
    //assert.NoError(t, err)

    err = reg.AddFileDescr(origFileDescr)
    assert.NoError(t, err)

    fileDescr, err := reg.GetFileDescr(fileId)
    assert.NoError(t, err)
    assert.Equal(t, origFileDescr, fileDescr)

    origMetaJSON, _ := json.MarshalIndent(origFileDescr, " ", "    ")
    metaJSON, _ := json.MarshalIndent(fileDescr, " ", "    ")
    assert.Equal(t, string(origMetaJSON), string(metaJSON))

    //err = reg.DeleteFileDescr(fileId)
    //assert.NoError(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}


func BenchmarkFileDescrInsertDelete(b *testing.B) {
    var err error

    var batchSize   int64 = 4
    var blockSize   int64 = 1024 * 1024

    var dataSize int64 = 1 //1024 * 1024 * 64
    data := make([]byte, dataSize)
    rand.Read(data)

    dbPath := "postgres://pgsql@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(b, err)

    err = reg.MigrateDB()
    assert.NoError(b, err)

    b.ResetTimer()

    const numRange int = 1024 * 1024 * 1024
    pBench := func(pb *testing.PB) {
        for pb.Next() {

            fileId, err := reg.GetNewFileId()
            assert.NoError(b, err)

            exists, err := reg.FileDescrExists(fileId)
            if exists {
                continue
            }
            assert.NoError(b, err)

            baseDir := b.TempDir()
            file := fsfile.NewFile(baseDir, fileId, batchSize, blockSize)

            err = file.Open()
            assert.NoError(b, err)

            err = file.Truncate()
            assert.NoError(b, err)

            reader := bytes.NewReader(data)

            written, err := file.Write(reader)
            assert.Error(b, err)
            assert.Equal(b, err, io.EOF)
            assert.Equal(b, dataSize, written)

            fileDescr, err := file.Meta()
            assert.NoError(b, err)

            err = reg.DeleteFileDescr(fileId)
            assert.NoError(b, err)

            err = reg.AddFileDescr(fileDescr)
            assert.NoError(b, err)

            _, err = reg.GetFileDescr(fileId)
            assert.NoError(b, err)

        }
    }
    b.SetParallelism(30)
    b.RunParallel(pBench)

    err = reg.CloseDB()
    assert.NoError(b, err)
}
