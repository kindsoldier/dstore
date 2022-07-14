/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bscont

import (
    "errors"
    "dstore/bstore/bsapi"
    "dstore/dsrpc"
    "dstore/dserr"
)

func (contr *Contr) SaveBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewSaveBlockParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    dataSize    := context.BinSize()
    blockReader := context.BinReader()

    fileId      := params.FileId
    batchId     := params.BatchId
    blockType   := params.BlockType
    blockId     := params.BlockId

    blockSize   := params.BlockSize

    err = contr.store.SaveBlock(fileId, batchId, blockType, blockId, blockSize, blockReader, dataSize)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    result := bsapi.NewSaveBlockResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) LoadBlockHandler(context *dsrpc.Context) error {
    var err error

    params := bsapi.NewLoadBlockParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    fileId      := params.FileId
    batchId     := params.BatchId
    blockType   := params.BlockType
    blockId     := params.BlockId

    blockWriter  := context.BinWriter()

    has, dataSize, err := contr.store.HasBlock(fileId, batchId, blockType, blockId)
    if err != nil {
        err = dserr.Err(err)
        context.SendError(err)
        return err
    }
    if !has {
        err = errors.New("block not exists")
        err = dserr.Err(err)
        context.SendError(err)
        return err
    }
    result := bsapi.NewLoadBlockResult()
    err = context.SendResult(result, dataSize)
    if err != nil {
        return dserr.Err(err)
    }
    err = contr.store.LoadBlock(fileId, batchId, blockType, blockId, blockWriter, dataSize)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) DeleteBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewDeleteBlockParams()

    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    fileId      := params.FileId
    batchId     := params.BatchId
    blockType   := params.BlockType
    blockId     := params.BlockId

    err = contr.store.DeleteBlock(fileId, batchId, blockType, blockId)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewDeleteBlockResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) ListBlocksHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewListBlocksParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    fileId := params.FileId

    blocks, err := contr.store.ListBlocks(fileId)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewListBlocksResult()
    result.Blocks = blocks
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
