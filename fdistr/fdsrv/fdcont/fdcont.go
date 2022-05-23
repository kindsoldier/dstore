/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "io"
    "ndstore/fdistr/fdapi"
    "ndstore/fdistr/fdsrv/fdrec"
    "ndstore/dsrpc"
)


type Contr struct {
    Store  *fdrec.Store
}

func NewContr() *Contr {
    return &Contr{}
}

const HelloMsg string = "hello"

func (contr *Contr) HelloHandler(context *dsrpc.Context) error {
    var err error
    params := fdapi.NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := fdapi.NewHelloResult()
    result.Message = HelloMsg
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) SaveHandler(context *dsrpc.Context) error {
    var err error
    params := fdapi.NewSaveParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }

    fileSize   := context.BinSize()
    fileReader := context.BinReader()

    filePath := params.FilePath
    err = contr.Store.SaveFile(filePath, fileReader, fileSize)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := fdapi.NewSaveResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) LoadHandler(context *dsrpc.Context) error {
    var err error
    params := fdapi.NewLoadParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    filePath := params.FilePath
    fileWriter := context.BinWriter()

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return err
    }

    fileSize, err := contr.Store.FileExists(filePath)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fdapi.NewLoadResult()
    err = context.SendResult(result, fileSize)
    if err != nil {
        return err
    }

    err = contr.Store.LoadFile(filePath, fileWriter)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) DeleteHandler(context *dsrpc.Context) error {
    var err error
    params := fdapi.NewDeleteParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }
    filePath   := params.FilePath
    err = contr.Store.DeleteFile(filePath)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fdapi.NewDeleteResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListHandler(context *dsrpc.Context) error {
    var err error
    params := fdapi.NewListParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    filePath   := params.DirPath

    files, err := contr.Store.ListFiles(filePath)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fdapi.NewListResult()
    result.Files = files
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
