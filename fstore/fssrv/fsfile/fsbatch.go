package fsfile

import (
    "errors"
    "io"
    "ndstore/dscom"
    "ndstore/dserr"
    "ndstore/fstore/fssrv/fsreg"
)

type Batch struct {
    reg         *fsreg.Reg
    baseDir     string
    fileId      int64
    batchId     int64
    batchSize   int64
    blocks      []*Block
}

func NewBatch(reg *fsreg.Reg, baseDir string, fileId, batchId, batchSize, blockSize int64) (*Batch, error) {
    var batch Batch
    var err error
    batch.reg       = reg
    batch.baseDir   = baseDir
    batch.fileId    = fileId
    batch.batchId   = batchId
    batch.batchSize = batchSize

    err = batch.addBatchDescr()
    if err != nil {
        return &batch, dserr.Err(err)
    }

    blockType := dscom.BTypeData
    batch.blocks = make([]*Block, batch.batchSize)
    for i := int64(0); i < batch.batchSize; i++ {
        blockId := i
        block, err := NewBlock(reg, baseDir, fileId, batchId, blockId, blockType, blockSize)
        if err != nil {
            return &batch, dserr.Err(err)
        }
        batch.blocks[i] = block
    }
    return &batch, dserr.Err(err)
}

func OpenBatch(reg *fsreg.Reg, baseDir string, fileId, batchId int64) (*Batch, error) {
    var err error
    var batch Batch
    exists, descr, err := reg.GetBatchDescr(fileId, batchId)
    if err != nil {
        return &batch, dserr.Err(err)
    }
    if !exists {
        err = errors.New("batch not exists")
        return &batch, dserr.Err(err)
    }

    batch.reg       = reg
    batch.baseDir   = baseDir

    batch.fileId    = descr.FileId
    batch.batchId   = descr.BatchId
    batch.batchSize = descr.BatchSize

    blockType := dscom.BTypeData
    batch.blocks = make([]*Block, batch.batchSize)
    for i := int64(0); i < batch.batchSize; i++ {
        blockId := i
        block, err := OpenBlock(reg, baseDir, fileId, batchId, blockId, blockType)
        if err != nil {
            return &batch, dserr.Err(err)
        }
        batch.blocks[i] = block
    }
    return &batch, dserr.Err(err)
}

func (batch *Batch) Write(reader io.Reader, need int64) (int64, error) {
    var err error
    var written int64
    for i := int64(0); i < batch.batchSize; i++ {
        if need < 1 {
            return written, err
        }
        blockWritten, err := batch.blocks[i].Write(reader, need)
        written += blockWritten
        if err == io.EOF {
            err = nil
            return written, dserr.Err(err)
        }
        if err != nil {
            return written, dserr.Err(err)
        }
        need -= blockWritten
    }
    return written, dserr.Err(err)
}


func (batch *Batch) Read(writer io.Writer) (int64, error) {
    var err error
    var read int64
    for i := int64(0); i < batch.batchSize; i++ {
        blockRead, err := batch.blocks[i].Read(writer)
        read += blockRead
        if err != nil {
            return read, dserr.Err(err)
        }
    }
    return read, dserr.Err(err)
}

func (batch *Batch) Close() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].Close()
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

//func (batch *Batch) Clean() error {
//    var err error
//    for i := int64(0); i < batch.batchSize; i++ {
//        err := batch.blocks[i].Clean()
//        if err != nil {
//            return dserr.Err(err)
//        }
//    }
//    return dserr.Err(err)
//}

func (batch *Batch) Erase() error {
    var err error
    for i := int64(0); i < batch.batchSize; i++ {
        err := batch.blocks[i].Erase()
        if err != nil {
            return dserr.Err(err)
        }
    }
    batch.blocks = make([]*Block, batch.batchSize)
    err = batch.reg.EraseBatchDescr(batch.fileId, batch.batchId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (batch *Batch) addBatchDescr() error {
    descr := batch.toDescr()
    return batch.reg.AddBatchDescr(descr)
}

func (batch *Batch) updateBatchDescr() error {
    descr := batch.toDescr()
    return batch.reg.UpdateBatchDescr(descr)
}

func (batch *Batch) toDescr() *dscom.BatchDescr {
    descr := dscom.NewBatchDescr()
    descr.FileId    = batch.fileId
    descr.BatchId   = batch.batchId
    descr.BatchSize = batch.batchSize
    return descr
}
