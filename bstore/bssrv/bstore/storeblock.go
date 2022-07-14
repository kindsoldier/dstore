/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bstore

import (
    "fmt"
    "io"
    "dstore/bstore/bssrv/bsblock"
    "dstore/dscomm/dserr"
    "dstore/dscomm/dsdescr"
)


func (store *Store) SaveBlock(fileId, batchId, blockType, blockId, blockSize int64, blockReader io.Reader, dataSize int64) error {
    var err error
    var has bool

    has, err = store.reg.HasBlock(fileId, batchId, blockType, blockId)
    if err != nil {
        return dserr.Err(err)
    }
    if has {
        block, err := bsblock.OpenBlock(store.reg, store.dataDir, fileId, batchId, blockType, blockId)
        if err != nil {
            return dserr.Err(err)
        }
        err = block.Clean()
        if err != nil {
            return dserr.Err(err)
        }
    }
    block, err := bsblock.NewBlock(store.reg, store.dataDir, fileId, batchId, blockType, blockId, blockSize)
    if err != nil {
        return dserr.Err(err)
    }
    wrSize, err := block.Write(blockReader, dataSize)
    if err == io.EOF {
        return dserr.Err(err)
    }
    if err != nil  {
        return dserr.Err(err)
    }
    if wrSize != dataSize {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) HasBlock(fileId, batchId, blockType, blockId int64) (bool, int64, error) {
    var err error
    var has bool
    var blockSize int64
    has, err = store.reg.HasBlock(fileId, batchId, blockType, blockId)
    if err != nil {
        return has, blockSize, dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("block %d,%d,%d,%d not exist", fileId, batchId, blockType, blockId)
        return has, blockSize, dserr.Err(err)
    }
    descr, err := store.reg.GetBlock(fileId, batchId, blockType, blockId)
    if err != nil {
        return has, blockSize, dserr.Err(err)
    }
    blockSize = descr.DataSize
    return has, blockSize, dserr.Err(err)
}

func (store *Store) LoadBlock(fileId, batchId, blockType, blockId int64, blockWriter io.Writer, dataSize int64) error {
    var err error
    has, err := store.reg.HasBlock(fileId, batchId, blockType, blockId)
    if err != nil {
        return dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("block %d,%d,%d,%d not exist", fileId, batchId, blockType, blockId)
        return dserr.Err(err)
    }
    block, err := bsblock.OpenBlock(store.reg, store.dataDir, fileId, batchId, blockType, blockId)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = block.Read(blockWriter, dataSize)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (store *Store) ListBlocks(fileId int64) ([]*dsdescr.Block, error) {
    var err error
    blocks := make([]*dsdescr.Block, 0)
    blocks, err = store.reg.ListBlocks(fileId)
    if err != nil {
        return blocks, dserr.Err(err)
    }
    return blocks, dserr.Err(err)
}

func (store *Store) DeleteBlock(fileId, batchId, blockType, blockId int64) error {
    var err error
    has, err := store.reg.HasBlock(fileId, batchId, blockType, blockId)
    if err != nil {
        return dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("block %d,%d,%d,%d not exist", fileId, batchId, blockType, blockId)
        return dserr.Err(err)
    }
    block, err := bsblock.OpenBlock(store.reg, store.dataDir, fileId, batchId, blockType, blockId)
    if err != nil {
        return dserr.Err(err)
    }
    err = block.Clean()
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
