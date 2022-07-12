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
)

type Block struct {
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


func NewBlock(baseDir string, blockId, batchId, fileId, blockSize int64) (*Block, error) {
    var err error
    var block Block
    block.baseDir   = baseDir

    block.blockId   = blockId
    block.batchId   = batchId
    block.fileId    = fileId

    block.blockSize = blockSize
    block.dataSize  = 0
    block.filePath  = newFilePath()
    block.createdAt = time.Now().Unix()
    block.updatedAt = block.createdAt
    return &block, err
}

func OpenBlock(baseDir string, descr *dsdescr.Block) (*Block, error) {
    var err error
    var block Block
    block.baseDir   = baseDir

    block.blockId   = descr.BlockId
    block.batchId   = descr.BatchId
    block.fileId    = descr.FileId

    block.blockSize = descr.BlockSize
    block.dataSize  = descr.DataSize
    block.filePath  = descr.FilePath
    block.createdAt = descr.CreatedAt
    block.updatedAt = descr.UpdatedAt
    return &block, err
}


func (block *Block) Write(reader io.Reader, dataSize int64) (int64, error) {
    var err error
    var wrSize int64
    remainSize := block.blockSize - block.dataSize
    if remainSize < dataSize {
        dataSize = remainSize
    }
    if remainSize < 1 || dataSize < 1 {
        return wrSize, err
    }
    newPath := newFilePath()
    writer, err := NewCrate(block.baseDir, newPath)
    if err != nil {
        err = fmt.Errorf("block write error: %s", err)
        return wrSize, err
    }
    var origin dsinter.Crate
    if block.dataSize > 0 {
        var wrSize int64
        reader, err := NewCrate(block.baseDir, block.filePath)
        if err != nil {
            err = fmt.Errorf("block recopy error: %s", err)
            return wrSize, err
        }
        wrSize, err = copyData(reader, writer, block.dataSize)
        if err != nil {
            err = fmt.Errorf("block recopy error: %s", err)
            return wrSize, err
        }
        if wrSize != dataSize {
            err = fmt.Errorf("block recopy only %d", wrSize)
        }
        origin = reader
    }
    wrSize, err = copyData(reader, writer, dataSize)
    if err != nil {
        writer.Clean()
        err = fmt.Errorf("block copy error: %s", err)
        return wrSize, err
    }
    if wrSize != dataSize {
        writer.Clean()
        err = fmt.Errorf("block copy only %d", wrSize)
        return wrSize, err
    }
    block.updatedAt = time.Now().Unix()
    block.filePath  = newPath
    block.dataSize += wrSize
    if origin != nil {
        origin.Clean()
    }
    return wrSize, err
}

func (block *Block) Read(writer io.Writer) (int64, error) {
    var err error
    var readSize int64

    reader, err := NewCrate(block.baseDir, block.filePath)
    if err != nil {
        err = fmt.Errorf("block read error: %s", err)
        return readSize, err
    }

    readSize, err = copyData(reader, writer, block.dataSize)
    if err != nil {
        err = fmt.Errorf("block recopy error: %s", err)
        return readSize, err
    }
    if readSize != block.dataSize {
        err = fmt.Errorf("block recopy only %d", readSize)
    }
    return readSize, err
}

func (block *Block) Descr() *dsdescr.Block {
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
    tank, err := NewCrate(block.baseDir, block.filePath)
    if err != nil {
        err = fmt.Errorf("block clean error: %s", err)
        return err
    }
    err = tank.Clean()
    if err != nil {
        err = fmt.Errorf("block clean error: %s", err)
        return err
    }
    block.dataSize = 0
    block.filePath = newFilePath()
    return err
}
