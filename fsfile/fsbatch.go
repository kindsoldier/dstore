/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import (
    //"fmt"
    "io"
    "time"

    //"dstore/dsdescr"
    //"dstore/dsinter"
)

type Pair struct {
    Id      int64
    Ref     interface{}
}

func NewPair(id int64, ref interface{}) *Pair {
    return &Pair{ Id: id, Ref: ref }
}


type Batch struct {
    baseDir     string
    batchSize   int64
    blockSize   int64
    createdAt   int64
    updatedAt   int64
    blocks      []*Block
    bb          []*Pair
}

func NewBatch(baseDir string, batchSize, blockSize int64) (*Batch, error) {
    var err error
    var batch Batch
    batch.baseDir   = baseDir

    batch.batchSize = batchSize
    batch.blockSize = blockSize
    batch.createdAt = time.Now().Unix()
    batch.updatedAt = batch.createdAt

    batch.blocks = make([]*Block, batchSize)
    for i := 0; int64(i) < batchSize; i++ {
        block, err := NewBlock(baseDir, blockSize)
        if err != nil {
            return &batch, err
        }
        batch.blocks[i] = block
    }
    return &batch, err
}

func (batch *Batch) Write(reader io.Reader, dataSize int64) (int64, error) {
    var err error
    var wrSize int64

    for i := 0; i < batch.countBlocks(); i++ {
        if dataSize < 1 {
            return dataSize, err
        }
        blockWrSize, err := batch.blocks[i].Write(reader, dataSize)
        wrSize += blockWrSize
        if err == io.EOF {
            err = nil
            return wrSize, err
        }
        if err != nil {
            return wrSize, err
        }
        dataSize -= wrSize
    }
    if wrSize > 0 {
        batch.updatedAt = time.Now().Unix()
    }
    return wrSize, err
}


func (batch *Batch) Read(writer io.Writer) (int64, error) {
    var err error
    var readSize int64
    for i := 0; i < batch.countBlocks(); i++ {
        blockReadSize, err := batch.blocks[i].Read(writer)
        readSize += blockReadSize
        if err != nil {
            return readSize, err
        }
    }
    return readSize, err
}

func (batch *Batch) countBlocks() int {
    return len(batch.blocks)
}
