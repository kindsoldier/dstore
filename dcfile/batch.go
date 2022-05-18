package dcfile


import (
    //"errors"
    //"fmt"
    "io"
    //"io/fs"
    "os"
    //"path/filepath"
)

type Batch struct {
    file        *os.File
    baseDir     string
    fileId      int64
    batchId     int64
    capacity    int64
    blocks      []*Block
}

func NewBatch(baseDir string, fileId, batchId, capacity, blockSize int64) *Batch {
    var batch Batch
    batch.baseDir   = baseDir
    batch.fileId    = fileId
    batch.batchId   = batchId
    batch.capacity  = capacity

    batch.blocks = make([]*Block, batch.capacity)
    for i := int64(0); i < batch.capacity; i++ {
        batch.blocks[i] = NewBlock(baseDir, fileId, batchId, i, blockSize)
    }
    return &batch
}

func (batch *Batch) Open() error {
    var err error
    for i := int64(0); i < batch.capacity; i++ {
        err = batch.blocks[i].Open()
        if err != nil {
            return err
        }
    }
    return err
}

func (batch *Batch) Write(reader io.Reader) (int64, error) {
    var err error
    var recorded int64
    for i := int64(0); i < batch.capacity; i++ {
        blockWritten, err := batch.blocks[i].Write(reader)
        if err != nil {
            return recorded, err
        }
        recorded += blockWritten
    }
    return recorded, err
}

func (batch *Batch) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    return read, err
}

func (batch *Batch) Close() error {
    var err error
    return err
}


func (batch *Batch) Truncate() error {
    var err error
    return err
}
