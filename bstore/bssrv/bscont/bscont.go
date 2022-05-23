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

const HelloMsg string = "hello"

type Contr struct {
    Store   *bsrec.Store
}

func NewContr() *Contr {
    return &Contr{}
}

func (contr *Contr) HelloHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := bsapi.NewHelloResult()
    result.Message = HelloMsg
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) SaveHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewSaveParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }

    blockSize   := context.BinSize()
    blockReader := context.BinReader()

    clusterId   := params.ClusterId
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    err = contr.Store.SaveBlock(clusterId, fileId, batchId, blockId,
                                                blockReader, blockSize)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := bsapi.NewSaveResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) LoadHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewLoadParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    clusterId   := params.ClusterId
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId

    blockWriter := context.BinWriter()

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return err
    }

    blockSize, err := contr.Store.BlockExists(clusterId, fileId, batchId, blockId)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewLoadResult()
    err = context.SendResult(result, blockSize)
    if err != nil {
        return err
    }

    err = contr.Store.LoadBlock(clusterId, fileId, batchId, blockId, blockWriter)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) DeleteHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewDeleteParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }
    clusterId   := params.ClusterId
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId

    err = contr.Store.DeleteBlock(clusterId, fileId, batchId, blockId)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewDeleteResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewListParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    clusterId   := params.ClusterId

    blocks, err := contr.Store.ListBlocks(clusterId)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewListResult()
    result.Blocks = blocks
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
