/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import (
    "io"
    "time"
    "dstore/dscomm/dsinter"
    "dstore/dscomm/dsdescr"
    "dstore/dscomm/dserr"
)

type Batch struct {
    reg         dsinter.FStoreReg
    baseDir     string
    batchId     int64
    fileId      int64
    batchSize   int64
    blockSize   int64
    createdAt   int64
    updatedAt   int64
    blocks      []*Block
}

func NewBatch(reg dsinter.FStoreReg, baseDir string, batchId, fileId, batchSize, blockSize int64) (*Batch, error) {
    var err error
    var batch Batch
    batch.baseDir   = baseDir
    batch.reg       = reg

    batch.batchId   = batchId
    batch.fileId    = fileId
    batch.batchSize = batchSize
    batch.blockSize = blockSize
    batch.createdAt = time.Now().Unix()
    batch.updatedAt = batch.createdAt

    batch.blocks = make([]*Block, batch.batchSize)
    for i := int64(0); i < batchSize; i++ {
        block, err := NewBlock(batch.reg, baseDir, i, batch.batchId, batch.fileId, blockSize)
        if err != nil {
            return &batch, dserr.Err(err)
        }
        batch.blocks[i] = block
    }
    descr := batch.toDescr()
    err = reg.PutBatch(descr)
    if err != nil {
        return &batch, dserr.Err(err)
    }
    return &batch, dserr.Err(err)
}

func OpenBatch(reg dsinter.FStoreReg, baseDir string, batchId, fileId int64) (*Batch, error) {
    var err error
    var batch Batch
    batch.baseDir   = baseDir
    batch.reg       = reg

    descr, err := reg.GetBatch(batchId, fileId)
    if err != nil {
        return &batch, dserr.Err(err)
    }

    batch.batchId   = descr.BatchId
    batch.fileId    = descr.FileId
    batch.batchSize = descr.BatchSize
    batch.blockSize = descr.BlockSize
    batch.createdAt = descr.CreatedAt
    batch.updatedAt = descr.UpdatedAt

    batch.blocks = make([]*Block, batch.batchSize)
    for i := int64(0); i < batch.batchSize; i++ {
        block, err := OpenBlock(reg, baseDir, i, batch.batchId, batch.fileId)
        if err != nil {
            return &batch, dserr.Err(err)
        }
        batch.blocks[i] = block
    }
    return &batch, dserr.Err(err)
}

func (batch *Batch) Write(reader io.Reader, reqSize int64) (int64, error) {
    var err error
    var wrSize int64

    for i := int64(0); i < batch.batchSize; i++ {
        if reqSize < 1 {
            return wrSize, dserr.Err(err)
        }
        blockWrSize, err := batch.blocks[i].Write(reader, reqSize)
        wrSize += blockWrSize
        if err == io.EOF {
            return wrSize, dserr.Err(err)
        }
        if err != nil {
            return wrSize, dserr.Err(err)
        }
        reqSize -= blockWrSize
    }

    if wrSize > 0 {
        batch.updatedAt = time.Now().Unix()
        descr := batch.toDescr()
        err = batch.reg.PutBatch(descr)
        if err != nil {
            return wrSize, dserr.Err(err)
        }
    }
    return wrSize, dserr.Err(err)
}

func (batch *Batch) Read(writer io.Writer, dataSize int64) (int64, error) {
    var err error
    var readSize int64

    if dataSize < 1 {
        return readSize, dserr.Err(err)
    }

    for i := int64(0); i < batch.batchSize; i++ {
        blockReadSize, err := batch.blocks[i].Read(writer, dataSize)
        readSize += blockReadSize
        dataSize -= blockReadSize
        if err != nil {
            return readSize, dserr.Err(err)
        }
    }
    return readSize, dserr.Err(err)
}

func (batch *Batch) Clean() error {
    var err error
    for i := batch.batchSize - 1; i >= 0; i-- {
        if batch.blocks[i] != nil {
            err := batch.blocks[i].Clean()
            if err != nil {
                return dserr.Err(err)
            }
            batch.blocks[i] = nil
        }
    }
    err = batch.reg.DeleteBatch(batch.batchId, batch.fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (batch *Batch) toDescr() *dsdescr.Batch {
    descr := dsdescr.NewBatch()
    descr.BatchId   = batch.batchId
    descr.FileId    = batch.fileId
    descr.BatchSize = batch.batchSize
    descr.BlockSize = batch.blockSize
    descr.CreatedAt = batch.createdAt
    descr.UpdatedAt = batch.updatedAt
    return descr
}
