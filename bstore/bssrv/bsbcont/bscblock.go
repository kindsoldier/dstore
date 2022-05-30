/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsbcont

import (
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsblock"
    "ndstore/dsrpc"
)

const GetHelloMsg string = "hello"

type Contr struct {
    store   *bsblock.Store
}

func NewContr(store *bsblock.Store) *Contr {
    return &Contr{ store: store }
}

func (contr *Contr) GetHelloHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewGetHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := bsapi.NewGetHelloResult()
    result.Message = GetHelloMsg
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) SaveBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewSaveBlockParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }

    blockSize   := context.BinSize()
    blockReader := context.BinReader()

    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    err = contr.store.SaveBlock(fileId, batchId, blockId, blockReader, blockSize)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := bsapi.NewSaveBlockResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) LoadBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewLoadBlockParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId

    blockWriter := context.BinWriter()

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return err
    }

    blockSize, err := contr.store.BlockExists(fileId, batchId, blockId)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewLoadBlockResult()
    err = context.SendResult(result, blockSize)
    if err != nil {
        return err
    }

    err = contr.store.LoadBlock(fileId, batchId, blockId, blockWriter)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) DeleteBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewDeleteBlockParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId

    err = contr.store.DeleteBlock(fileId, batchId, blockId)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewDeleteBlockResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListBlocksHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewListBlocksParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    blocks, err := contr.store.ListBlocks()
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewListBlocksResult()
    result.Blocks = blocks
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
