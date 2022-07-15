/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fscont

import (
    "errors"
    "dstore/fstore/fsapi"
    "dstore/dscomm/dsrpc"
    "dstore/dscomm/dserr"
)

func (contr *Contr) SaveFileHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewSaveFileParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    fileSize    := context.BinSize()
    fileReader  := context.BinReader()
    login    := string(context.AuthIdent())

    filePath := params.FilePath
    descr, err := contr.store.SaveFile(login, filePath, fileReader, fileSize)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    result := fsapi.NewSaveFileResult()
    result.File = descr
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) LoadFileHandler(context *dsrpc.Context) error {
    var err error

    params := fsapi.NewLoadFileParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    filePath    := params.FilePath
    fileWriter  := context.BinWriter()
    login := string(context.AuthIdent())

    has, descr, err := contr.store.HasFile(login, filePath)
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
    if descr == nil {
        err = errors.New("file descr is nil")
        err = dserr.Err(err)
        context.SendError(err)
        return err
    }

    fileSize := descr.DataSize
    result := fsapi.NewLoadFileResult()
    result.File = descr
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
    params := fsapi.NewDeleteFileParams()

    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    filePath    := params.FilePath
    login    := string(context.AuthIdent())

    descr, err:= contr.store.DeleteFile(login, filePath)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewDeleteFileResult()
    result.File = descr
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) ListFilesHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewListFilesParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    pattern     := params.Pattern
    regular     := params.Regular
    gPattern    := params.GPattern

    login   := string(context.AuthIdent())

    files, err := contr.store.ListFiles(login, pattern, regular, gPattern)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewListFilesResult()
    result.Files = files
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) FileStatsHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewFileStatsParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    pattern     := params.Pattern
    regular     := params.Regular
    gPattern    := params.GPattern

    login   := string(context.AuthIdent())

    count, usage, err := contr.store.FileStats(login, pattern, regular, gPattern)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewFileStatsResult()
    result.Usage = usage
    result.Count = count
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
