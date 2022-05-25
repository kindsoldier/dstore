package fsreg

import (
    "encoding/json"
    "fmt"
    "bytes"
    "io"
    //"path/filepath"
    "testing"
    "math/rand"

    "github.com/stretchr/testify/assert"

    //"ndstore/dscom"
    "ndstore/fstore/fssrv/fsfile"
)

func TestInsertSelectDelete(t *testing.T) {
    var err error

    var fileId      int64 = 1
    var batchSize   int64 = 5
    var blockSize   int64 = 1024

    baseDir := t.TempDir()

    file := fsfile.NewFile(baseDir, fileId, batchSize, blockSize)

    err = file.Open()
    assert.NoError(t, err)

    err = file.Truncate()
    assert.NoError(t, err)

    var dataSize int64 = 2 * (batchSize * blockSize) + blockSize + 2
    data := make([]byte, dataSize)
    rand.Read(data)
    reader := bytes.NewReader(data)

    written, err := file.Write(reader)
    assert.Error(t, err)
    assert.Equal(t, err, io.EOF)
    assert.Equal(t, dataSize, written)

    fileDescr, err := file.Meta()
    assert.NoError(t, err)


    err = file.Close()
    assert.NoError(t, err)

    dbPath := "./files.db"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    err = reg.AddFileDescr(fileDescr)
    assert.NoError(t, err)

    fileDescr, err = reg.GetFileDescr(fileId)
    assert.NoError(t, err)

    metaJSON, _ := json.MarshalIndent(fileDescr, " ", "    ")
    fmt.Println(string(metaJSON))

    err = reg.CloseDB()
    assert.NoError(t, err)

}
