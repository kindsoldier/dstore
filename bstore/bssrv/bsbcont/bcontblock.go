/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsbcont

import (
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/dscom"
    "ndstore/dsrpc"
    "ndstore/dserr"
)

func (contr *Contr) SaveBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewSaveBlockParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    binSize     := context.BinSize()
    blockReader := context.BinReader()

    descr := dscom.NewBlockDescr()

    descr.FileId      = params.FileId
    descr.BatchId     = params.BatchId
    descr.BlockId     = params.BlockId
    descr.BlockType   = params.BlockType
    descr.BlockVer    = params.BlockVer

    descr.BlockSize   = params.BlockSize
    descr.DataSize    = params.DataSize
    descr.HashAlg     = params.HashAlg
    descr.HashInit    = params.HashInit
    descr.HashSum     = params.HashSum


    err = contr.store.SaveBlock(descr, blockReader, binSize)
    if err != nil {
        context.SendError(dserr.Err(err))
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
        return dserr.Err(err)
    }

    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType
    blockVer    := params.BlockVer

    blockWriter := context.BinWriter()

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    _, blockVer, dataSize, err := contr.store.GetBlockParams(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewLoadBlockResult()
    err = context.SendResult(result, dataSize)
    if err != nil {
        return dserr.Err(err)
    }

    err = contr.store.LoadBlock(fileId, batchId, blockId, blockType, blockVer, blockWriter)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) BlockExistsHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewBlockExistsParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType
    blockVer    := params.BlockVer

    exists, _, _, err := contr.store.GetBlockParams(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewBlockExistsResult()
    result.Exists = exists
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) CheckBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewCheckBlockParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType
    blockVer    := params.BlockVer

    correct, err := contr.store.CheckBlock(fileId, batchId, blockId, blockType, blockVer)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewCheckBlockResult()
    result.Correct = correct
    err = context.SendResult(result, 0)
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
        return dserr.Err(err)
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType
    blockVer    := params.BlockVer

    err = contr.store.DeleteBlock(fileId, batchId, blockId, blockType, blockVer)
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

func (contr *Contr) LinkBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewLinkBlockParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType
    oldBlockVer    := params.OldBlockVer
    newBlockVer    := params.NewBlockVer

    err = contr.store.LinkBlock(fileId, batchId, blockId, blockType, oldBlockVer, newBlockVer)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewLinkBlockResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (contr *Contr) EraseAllHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewEraseAllParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }

    err = contr.store.EraseAll()
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewEraseAllResult()
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
        return dserr.Err(err)
    }

    blocks, err := contr.store.ListBlocks()
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
