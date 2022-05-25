/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bscont

import (
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/dsrpc"
    "ndstore/bstore/bssrv/bsrec"
)

const GetHelloMsg string = "hello"

type Contr struct {
    Store   *bsrec.Store
}

func NewContr() *Contr {
    return &Contr{}
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
    err = contr.Store.SaveBlock(fileId, batchId, blockId, blockReader, blockSize)
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

    blockSize, err := contr.Store.BlockExists(fileId, batchId, blockId)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewLoadBlockResult()
    err = context.SendResult(result, blockSize)
    if err != nil {
        return err
    }

    err = contr.Store.LoadBlock(fileId, batchId, blockId, blockWriter)
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

    err = contr.Store.DeleteBlock(fileId, batchId, blockId)
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

    blocks, err := contr.Store.ListBlocks()
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
