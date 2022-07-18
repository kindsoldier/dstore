/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fscont

import (
    "dstore/fstore/fsapi"
    "dstore/dscomm/dsrpc"
    "dstore/dscomm/dserr"

    "dstore/dscomm/dsdescr"
)

func (contr *Contr) AddBStoreHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewAddBStoreParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    descr := dsdescr.NewBStore()
    descr.Address = params.Address
    descr.Port    = params.Port
    descr.Login   = params.Login
    descr.Pass    = params.Pass
    descr.State   = params.State

    authLogin := string(context.AuthIdent())
    err = contr.store.AddBStore(authLogin, descr)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    result := fsapi.NewAddBStoreResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) UpdateBStoreHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewUpdateBStoreParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    descr := dsdescr.NewBStore()
    descr.Address = params.Address
    descr.Port    = params.Port
    descr.Login   = params.Login
    descr.Pass    = params.Pass
    descr.State   = params.State

    authLogin := string(context.AuthIdent())
    err = contr.store.UpdateBStore(authLogin, descr)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    result := fsapi.NewUpdateBStoreResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) DeleteBStoreHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewDeleteBStoreParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    address := params.Address
    port    := params.Port
    authLogin := string(context.AuthIdent())
    err = contr.store.DeleteBStore(authLogin, address, port)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewDeleteBStoreResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) ListBStoresHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewListBStoresParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    regular := params.Regular
    authLogin := string(context.AuthIdent())
    bstores, err := contr.store.ListBStores(authLogin, regular)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewListBStoresResult()
    result.BStores = bstores
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) CheckBStoreHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewCheckBStoreParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    address     := params.Address
    port        := params.Port
    login       := params.Login
    pass        := params.Pass

    authLogin    := string(context.AuthIdent())
    ok, err := contr.store.CheckBStore(authLogin, address, port, login, pass)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewCheckBStoreResult()
    result.Match = ok
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
