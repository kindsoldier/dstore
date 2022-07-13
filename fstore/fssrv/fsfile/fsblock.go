/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import (
    "fmt"
    "io"
    "time"

    "dstore/dsdescr"
    "dstore/dsinter"
    "dstore/dserr"
)

type Block struct {
    reg         dsinter.StoreReg
    baseDir     string
    filePath    string
    blockId     int64
    batchId     int64
    fileId      int64
    blockSize   int64
    dataSize    int64
    createdAt   int64
    updatedAt   int64
}

func NewBlock(reg dsinter.StoreReg, baseDir string, blockId, batchId, fileId, blockSize int64) (*Block, error) {
    var err error
    var block Block
    block.baseDir   = baseDir
    block.reg       = reg

    block.blockId   = blockId
    block.batchId   = batchId
    block.fileId    = fileId

    block.blockSize = blockSize
    block.dataSize  = 0
    block.filePath  = newFilePath()
    block.createdAt = time.Now().Unix()
    block.updatedAt = block.createdAt

    descr := block.toDescr()
    err = reg.PutBlock(descr)
    if err != nil {
        return &block, dserr.Err(err)
    }
    return &block, dserr.Err(err)
}

func OpenBlock(reg dsinter.StoreReg, baseDir string, blockId, batchId, fileId int64) (*Block, error) {
    var err error
    var block Block
    block.baseDir   = baseDir
    block.reg       = reg

    descr, err := reg.GetBlock(blockId, batchId, fileId)
    if err != nil {
        return &block, dserr.Err(err)
    }

    block.blockId   = descr.BlockId
    block.batchId   = descr.BatchId
    block.fileId    = descr.FileId

    block.blockSize = descr.BlockSize
    block.dataSize  = descr.DataSize
    block.filePath  = descr.FilePath
    block.createdAt = descr.CreatedAt
    block.updatedAt = descr.UpdatedAt
    return &block, dserr.Err(err)
}


func (block *Block) Write(reader io.Reader, dataSize int64) (int64, error) {
    var err error
    var wrSize int64
    remainSize := block.blockSize - block.dataSize
    if remainSize < dataSize {
        dataSize = remainSize
    }
    if remainSize < 1 || dataSize < 1 {
        return wrSize, dserr.Err(err)
    }
    newPath := newFilePath()
    writer, err := NewCrate(block.baseDir, newPath)
    if err != nil {
        err = fmt.Errorf("block write error: %s", err)
        return wrSize, dserr.Err(err)
    }
    var origin dsinter.Crate
    if block.dataSize > 0 {
        var wrSize int64
        reader, err := NewCrate(block.baseDir, block.filePath)
        if err != nil {
            err = fmt.Errorf("block recopy error: %s", err)
            return wrSize, dserr.Err(err)
        }
        wrSize, err = copyData(reader, writer, block.dataSize)
        if err != nil {
            err = fmt.Errorf("block recopy error: %s", err)
            return wrSize, dserr.Err(err)
        }
        if wrSize != dataSize {
            err = fmt.Errorf("block recopy only %d", wrSize)
            return wrSize, dserr.Err(err)
        }
        origin = reader
    }
    wrSize, err = copyData(reader, writer, dataSize)
    if err != nil {
        writer.Clean()
        err = fmt.Errorf("block copy error: %s", err)
        return wrSize, dserr.Err(err)
    }
    if wrSize != dataSize {
        writer.Clean()
        err = fmt.Errorf("block copy only %d", wrSize)
        return wrSize, dserr.Err(err)
    }
    block.updatedAt = time.Now().Unix()
    block.filePath  = newPath
    block.dataSize += wrSize
    if origin != nil {
        origin.Clean()
    }

    descr := block.toDescr()
    err = block.reg.PutBlock(descr)
    if err != nil {
        return wrSize, dserr.Err(err)
    }
    return wrSize, dserr.Err(err)
}

func (block *Block) Read(writer io.Writer) (int64, error) {
    var err error
    var readSize int64

    descr, err := block.reg.GetBlock(block.blockId, block.batchId, block.fileId)
    if err != nil {
        return readSize, dserr.Err(err)
    }

    block.blockSize = descr.BlockSize
    block.dataSize  = descr.DataSize
    block.filePath  = descr.FilePath
    block.createdAt = descr.CreatedAt
    block.updatedAt = descr.UpdatedAt

    reader, err := NewCrate(block.baseDir, block.filePath)
    if err != nil {
        err = fmt.Errorf("block read error: %s", err)
        return readSize, dserr.Err(err)
    }

    readSize, err = copyData(reader, writer, block.dataSize)
    if err != nil {
        err = fmt.Errorf("block recopy error: %s", err)
        return readSize, dserr.Err(err)
    }
    if readSize != block.dataSize {
        err = fmt.Errorf("block recopy only %d", readSize)
    }
    return readSize, dserr.Err(err)
}

func (block *Block) toDescr() *dsdescr.Block {
    descr := dsdescr.NewBlock()
    descr.BlockId   = block.blockId
    descr.BatchId   = block.batchId
    descr.FileId    = block.fileId
    descr.BlockSize = block.blockSize
    descr.DataSize  = block.dataSize
    descr.FilePath  = block.filePath
    descr.CreatedAt = block.createdAt
    descr.UpdatedAt = block.updatedAt
    return descr
}

func (block *Block) Clean() error {
    var err error
    crate, err := NewCrate(block.baseDir, block.filePath)
    if err != nil {
        err = fmt.Errorf("block clean error: %s", err)
        return dserr.Err(err)
    }
    err = crate.Clean()
    if err != nil {
        err = fmt.Errorf("block clean error: %s", err)
        return dserr.Err(err)
    }
    block.dataSize = 0
    block.filePath = newFilePath()

    err = block.reg.DeleteBlock(block.blockId, block.batchId, block.fileId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
