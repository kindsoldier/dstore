/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package ndcontr

import (
    "io"
    "dcstore/dcnode/ndapi"
    "dcstore/dcrpc"
    "dcstore/dcnode/ndsrv/ndstore"
)

type Contr struct {
    Store   *ndstore.Store
}

func NewContr() *Contr {
    return &Contr{}
}

func (contr *Contr) HelloHandler(context *dcrpc.Context) error {
    var err error
    params := ndapi.NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := ndapi.NewHelloResult()
    result.Message = "hello!"
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) SaveHandler(context *dcrpc.Context) error {
    var err error
    params := ndapi.NewSaveParams()

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

    result := ndapi.NewSaveResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) LoadHandler(context *dcrpc.Context) error {
    var err error
    params := ndapi.NewLoadParams()
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
    result := ndapi.NewLoadResult()
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

func (contr *Contr) DeleteHandler(context *dcrpc.Context) error {
    var err error
    params := ndapi.NewDeleteParams()

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
    result := ndapi.NewDeleteResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
