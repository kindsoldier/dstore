/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import (
    "io"
    "time"
    "dstore/dsinter"
    "dstore/dsdescr"
    "dstore/dslog"
    "dstore/dserr"
)

type Batch struct {
    reg         dsinter.StoreReg
    baseDir     string
    batchId     int64
    fileId      int64
    batchSize   int64
    blockSize   int64
    createdAt   int64
    updatedAt   int64
    blocks      []*Block
}

func NewBatch(baseDir string, reg dsinter.StoreReg, batchId, fileId, batchSize, blockSize int64) (*Batch, error) {
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
        block, err := NewBlock(baseDir, batch.reg, i, batch.batchId, batch.fileId, blockSize)
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

func OpenBatch(baseDir string, reg dsinter.StoreReg, batchId, fileId int64) (*Batch, error) {
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
        block, err := OpenBlock(baseDir, reg, i, batch.batchId, batch.fileId)
        if err != nil {
            return &batch, dserr.Err(err)
        }
        batch.blocks[i] = block
    }
    return &batch, dserr.Err(err)
}

func (batch *Batch) Write(reader io.Reader, dataSize int64) (int64, error) {
    var err error
    var wrSize int64

    for i := 0; i < batch.countBlocks(); i++ {
        if dataSize < 1 {
            return wrSize, dserr.Err(err)
        }
        blockWrSize, err := batch.blocks[i].Write(reader, dataSize)
        dslog.LogDebugf("written to block: %d", blockWrSize)
        wrSize += blockWrSize
        if err == io.EOF {
            err = nil
            return wrSize, dserr.Err(err)
        }
        if err != nil {
            return wrSize, dserr.Err(err)
        }
        dataSize -= blockWrSize
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

func (batch *Batch) Read(writer io.Writer) (int64, error) {
    var err error
    var readSize int64
    for i := 0; i < batch.countBlocks(); i++ {
        blockReadSize, err := batch.blocks[i].Read(writer)
        readSize += blockReadSize
        if err != nil {
            return readSize, dserr.Err(err)
        }
    }
    return readSize, dserr.Err(err)
}

func (batch *Batch) Clean() error {
    var err error
    for i := 0; i < batch.countBlocks(); i++ {
        err := batch.blocks[i].Clean()
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (batch *Batch) countBlocks() int {
    return len(batch.blocks)
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
