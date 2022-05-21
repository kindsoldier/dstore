package dcfile

import (
    "bytes"
    "math/rand"
    "io"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFileSmallWriteRead(t *testing.T) {
    var err error
    var fileId      int64 = 1
    var batchSize   int64 = 5
    var blockSize   int64 = 1024 * 1024

    baseDir := "./"
    //baseDir := t.TempDir()

    file := NewFile(baseDir, fileId, batchSize, blockSize)

    err = file.Open()
    assert.NoError(t, err)

    err = file.Truncate()
    assert.NoError(t, err)

    var dataSize int64 = 1024 * 1014 * 1024 //2 * (batchSize * blockSize) + blockSize + 2
    data := make([]byte, dataSize)
    rand.Read(data)
    reader := bytes.NewReader(data)

    written, err := file.Write(reader)
    assert.Error(t, err)
    assert.Equal(t, err, io.EOF)
    assert.Equal(t, dataSize, written)

    //err = file.Close()
    //assert.NoError(t, err)

    //file = NewFile(baseDir, fileId, batchSize, blockSize)

    //err = file.Open()
    //assert.NoError(t, err)

    //fileSize, _ := file.Size()
    //assert.Equal(t, written, fileSize)

    //writer := bytes.NewBuffer(make([]byte, 0))
    //read, err := file.Read(writer)
    //assert.Equal(t, dataSize, read)
    //assert.Equal(t, written, read)

    //err = file.Close()
    //assert.NoError(t, err)

    //assert.Equal(t, data[0:written], writer.Bytes())

    //file.Purge()
}
