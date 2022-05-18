package dcfile

import (
    "bytes"
    "math/rand"
    "io"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestBlockWriteReadLess(t *testing.T) {
    var err error
    var fileId  int64 = 1
    var batchId int64 = 2
    var blockId int64 = 3
    var bCap int64 = 1024

    baseDir := "./"
    block := NewBlock(baseDir, fileId, batchId, blockId, bCap)
    err = block.Open()
    assert.NoError(t, err)

    err = block.Truncate()
    assert.NoError(t, err)

    var dataSize int64 = 1024 - 2
    data := make([]byte, dataSize)
    rand.Read(data)

    reader := bytes.NewReader(data)
    written, err := block.Write(reader)
    assert.Equal(t, err, io.EOF)
    assert.Equal(t, dataSize, written)
    err = block.Close()

    block = NewBlock(baseDir, fileId, batchId, blockId, bCap)
    err = block.Open()
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))
    read, err := block.Read(writer)
    assert.Equal(t, dataSize, read)
    assert.Equal(t, written, read)

    err = block.Close()
    assert.NoError(t, err)

    assert.Equal(t, data[0:written], writer.Bytes())
}

func TestBlockWriteReadMore(t *testing.T) {
    var err error
    var fileId  int64 = 1
    var batchId int64 = 2
    var blockId int64 = 3
    var bCap int64 = 1024

    baseDir := "./"
    block := NewBlock(baseDir, fileId, batchId, blockId, bCap)
    err = block.Open()
    assert.NoError(t, err)

    err = block.Truncate()
    assert.NoError(t, err)

    var dataSize int64 = 1024 + 2
    data := make([]byte, dataSize)
    rand.Read(data)

    reader := bytes.NewReader(data)
    written, err := block.Write(reader)
    assert.NotEqual(t, err, io.EOF)
    //assert.Equal(t, dataSize, written)
    err = block.Close()

    block = NewBlock(baseDir, fileId, batchId, blockId, bCap)
    err = block.Open()
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))
    read, err := block.Read(writer)
    assert.Equal(t, written, read)

    err = block.Close()
    assert.NoError(t, err)

    assert.Equal(t, data[0:written], writer.Bytes())
}
