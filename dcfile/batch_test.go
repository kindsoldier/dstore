package dcfile

import (
    "bytes"
    "math/rand"
    "io"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestBatchSmallWriteRead(t *testing.T) {
    var err error
    var fileId      int64 = 1
    var batchId     int64 = 2
    var capa        int64 = 5
    var blockSize   int64 = 1024

    baseDir := t.TempDir()

    batch := NewBatch(baseDir, fileId, batchId, capa, blockSize)

    err = batch.Open()
    assert.NoError(t, err)

    err = batch.Truncate()
    assert.NoError(t, err)

    var dataSize int64 = blockSize * capa - 2
    data := make([]byte, dataSize)
    rand.Read(data)
    reader := bytes.NewReader(data)

    written, err := batch.Write(reader)
    assert.Error(t, err)
    assert.Equal(t, err, io.EOF)
    assert.Equal(t, dataSize, written)

    err = batch.Close()
    assert.NoError(t, err)

    batch = NewBatch(baseDir, fileId, batchId, capa, blockSize)

    err = batch.Open()
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))
    read, err := batch.Read(writer)
    assert.Equal(t, dataSize, read)
    assert.Equal(t, written, read)

    err = batch.Close()
    assert.NoError(t, err)

    assert.Equal(t, data[0:written], writer.Bytes())

    batch.Purge()
}
