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
    err = context.SendResult(result)
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

    blockSize := context.BinSize()
    blockReader := context.BinReader()

    contr.Store.SaveBlock(blockReader, blockSize)

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := ndapi.NewSaveResult()
    err = context.SendResult(result)
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

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := ndapi.NewLoadResult()
    err = context.SendResult(result)
    if err != nil {
        return err
    }
    return err
}
