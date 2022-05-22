/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fscont

import (
    "io"
    "ndstore/fstore/fsapi"
    "ndstore/dsrpc"
    "ndstore/fstore/fssrv/fsrec"
)

type Contr struct {
    Store   *fsrec.Store
}

func NewContr() *Contr {
    return &Contr{}
}

func (contr *Contr) HelloHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := fsapi.NewHelloResult()
    result.Message = "hello!"
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) SaveHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewSaveParams()

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

    result := fsapi.NewSaveResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) LoadHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewLoadParams()
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
    result := fsapi.NewLoadResult()
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
    params := fsapi.NewDeleteParams()

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
    result := fsapi.NewDeleteResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewListParams()
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
    result := fsapi.NewListResult()
    result.Blocks = blocks
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
