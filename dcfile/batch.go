package dcfile


import (
    "io"
)


type BatchMeta struct {
    Blocks      []*BlockMeta    `json:"blocks"`
}

func NewBatchMeta() *BatchMeta {
    var batchMeta BatchMeta
    batchMeta.Blocks = make([]*BlockMeta, 0)
    return &batchMeta
}


type Batch struct {
    baseDir     string
    fileId      int64
    batchId     int64
    batchSize   int64
    blocks      []*Block
}


func NewBatch(baseDir string, fileId, batchId, batchSize, blockSize int64) *Batch {
    var batch Batch
    batch.baseDir   = baseDir
    batch.fileId    = fileId
    batch.batchId   = batchId
    batch.batchSize  = batchSize

    batch.blocks = make([]*Block, batch.batchSize)
    for i := int64(0); i < batch.batchSize; i++ {
        batch.blocks[i] = NewBlock(baseDir, fileId, batchId, i, blockSize)
    }
    return &batch
}

func (batch *Batch) Meta() *BatchMeta {
    batchMeta := NewBatchMeta()
    for i := range batch.blocks {
        blockMeta := batch.blocks[i].Meta()
        batchMeta.Blocks = append(batchMeta.Blocks, blockMeta)
    }
    return batchMeta
}


func (batch *Batch) Open() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err = batch.blocks[i].Open()
        if err != nil {
            return err
        }
    }
    return err
}

func (batch *Batch) Write(reader io.Reader) (int64, error) {
    var err error
    var written int64
    for i := int64(0); i < batch.batchSize; i++ {
        blockWritten, err := batch.blocks[i].Write(reader)
        written += blockWritten
        if err != nil {
            return written, err
        }
    }
    return written, err
}

func (batch *Batch) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    for i := int64(0); i < batch.batchSize; i++ {
        blockRead, err := batch.blocks[i].Read(writer)
        read += blockRead
        if err != nil {
            return read, err
        }
    }
    return read, err
}

func (batch *Batch) Close() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].Close()
        if err != nil {
            return err
        }
    }
    return err
}

func (batch *Batch) Truncate() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].Truncate()
        if err != nil {
            return err
        }
    }
    return err
}

func (batch *Batch) Purge() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].Purge()
        if err != nil {
            return err
        }
    }
    return err
}

func (batch *Batch) Size() (int64, error) {
    var err error
    var size int64
    for i := int64(0); i < batch.batchSize; i++ {
        blockSize, err := batch.blocks[i].Size()
        size += blockSize
        if err != nil {
            return size, err
        }
    }
    return size, err
}

func (batch *Batch) ToBegin() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].ToBegin()
        if err != nil {
            return err
        }
    }
    return err
}

func (batch *Batch) ToEnd() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].ToEnd()
        if err != nil {
            return err
        }
    }
    return err
}
