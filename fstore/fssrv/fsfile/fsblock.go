/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import (
    "fmt"
    "io"
    "time"

    "dstore/dscomm/dsdescr"
    "dstore/dscomm/dserr"
)

type Block struct {
    baseDir     string
    filePath    string

    fileId      int64
    batchId     int64
    blockType   int64
    blockId     int64

    blockSize   int64
    dataSize    int64
    createdAt   int64
    updatedAt   int64
}

func NewBlock(baseDir string, fileId, batchId, blockType, blockId, blockSize int64) (*Block, error) {
    var err error
    var block Block
    block.baseDir   = baseDir

    block.fileId    = fileId
    block.batchId   = batchId
    block.blockType = blockType
    block.blockId   = blockId

    block.blockSize = blockSize
    block.dataSize  = 0
    block.filePath  = newFilePath()
    block.createdAt = time.Now().Unix()
    block.updatedAt = block.createdAt

    return &block, dserr.Err(err)
}

func OpenBlock(baseDir string, descr *dsdescr.Block) (*Block, error) {
    var err error
    var block Block
    block.baseDir   = baseDir

    block.fileId    = descr.FileId
    block.batchId   = descr.BatchId
    block.blockType = descr.BlockType
    block.blockId   = descr.BlockId

    block.blockSize = descr.BlockSize
    block.dataSize  = descr.DataSize
    block.filePath  = descr.FilePath

    block.createdAt = descr.CreatedAt
    block.updatedAt = descr.UpdatedAt
    return &block, dserr.Err(err)
}

func (block *Block) Write(reader io.Reader, dataSize int64) (int64, bool, error) {
    var err error
    var wrSize int64
    var eof bool
    remainSize := block.blockSize - block.dataSize
    if remainSize < dataSize {
        dataSize = remainSize
    }

    if remainSize < 1 || dataSize < 1 {
        return wrSize, eof, dserr.Err(err)
    }
    newPath := newFilePath()
    newPath = fmt.Sprintf("%s--%05d-%04d-%03d", newPath, block.fileId, block.batchId, block.blockId)

    writer, err := OpenCrate(block.baseDir, newPath, WRONLY)
    defer writer.Close()
    if err != nil {
        err = fmt.Errorf("block write error: %s", err)
        return wrSize, eof, dserr.Err(err)
    }

    if block.dataSize > 0 {
        var wrSize int64
        pReader, err := OpenCrate(block.baseDir, block.filePath, RDONLY)
        defer pReader.Close()
        if err != nil {
            err = fmt.Errorf("block recopy error: %s", err)
            return wrSize, eof, dserr.Err(err)
        }
        wrSize, eof, err = copyData(pReader, writer, block.dataSize)
        block.updatedAt = time.Now().Unix()
        block.filePath  = newPath
        block.dataSize += wrSize
        pReader.Clean()
        if err != nil {
            err = fmt.Errorf("block recopy error: %s", err)
            return wrSize, eof, dserr.Err(err)
        }
    }

    wrSize, eof, err = copyData(reader, writer, dataSize)
    if err == io.EOF {
        eof = true
        err = nil
    }
    block.updatedAt = time.Now().Unix()
    block.filePath  = newPath
    block.dataSize += wrSize
    if err != nil {
        err = fmt.Errorf("block copy error: %s", err)
        return wrSize, eof, dserr.Err(err)
    }
    return wrSize, eof, dserr.Err(err)
}

func (block *Block) Read(writer io.Writer, dataSize int64) (int64, error) {
    var err error
    var readSize int64

    if dataSize < 1 {
        return readSize, dserr.Err(err)
    }

    reader, err := OpenCrate(block.baseDir, block.filePath, RDONLY)
    defer reader.Close()
    if err != nil {
        err = fmt.Errorf("block read error: %s", err)
        return readSize, dserr.Err(err)
    }

    readSize, _, err = copyData(reader, writer, block.dataSize)
    if err != nil {
        err = fmt.Errorf("block recopy error: %s", err)
        return readSize, dserr.Err(err)
    }
    if readSize != block.dataSize {
        err = fmt.Errorf("block recopy only %d", readSize)
    }
    return readSize, dserr.Err(err)
}

func (block *Block) Descr() *dsdescr.Block {
    descr := dsdescr.NewBlock()
    descr.FileId    = block.fileId
    descr.BatchId   = block.batchId
    descr.BlockType = block.blockType
    descr.BlockId   = block.blockId

    descr.BlockSize = block.blockSize
    descr.DataSize  = block.dataSize
    descr.FilePath  = block.filePath

    descr.CreatedAt = block.createdAt
    descr.UpdatedAt = block.updatedAt
    return descr
}

func (block *Block) Clean() error {
    var err error
    crate, err := OpenCrate(block.baseDir, block.filePath, WRONLY)
    defer crate.Close()
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

    return dserr.Err(err)
}
