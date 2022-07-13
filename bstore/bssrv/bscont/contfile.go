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

func (contr *Contr) SaveFileHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewSaveFileParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    fileSize    := context.BinSize()
    fileReader  := context.BinReader()
    login    := string(context.AuthIdent())

    filePath := params.FilePath
    err = contr.store.SaveFile(login, filePath, fileReader, fileSize)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    result := bsapi.NewSaveFileResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) LoadFileHandler(context *dsrpc.Context) error {
    var err error

    params := bsapi.NewLoadFileParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    filePath    := params.FilePath
    fileWriter  := context.BinWriter()
    login := string(context.AuthIdent())

    has, fileSize, err := contr.store.HasFile(login, filePath)
    if err != nil {
        err = dserr.Err(err)
        context.SendError(err)
        return err
    }
    if !has {
        err = errors.New("file not exists")
        err = dserr.Err(err)
        context.SendError(err)
        return err
    }
    result := bsapi.NewLoadFileResult()
    err = context.SendResult(result, fileSize)
    if err != nil {
        return dserr.Err(err)
    }
    err = contr.store.LoadFile(login, filePath, fileWriter)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) DeleteFileHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewDeleteFileParams()

    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    filePath    := params.FilePath
    login    := string(context.AuthIdent())

    err = contr.store.DeleteFile(login, filePath)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewDeleteFileResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) ListFilesHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewListFilesParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    dirPath := params.DirPath
    login   := string(context.AuthIdent())

    files, err := contr.store.ListFiles(login, dirPath)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewListFilesResult()
    result.Files = files
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
