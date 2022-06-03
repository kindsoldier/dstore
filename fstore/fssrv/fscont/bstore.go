/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "ndstore/fstore/fsapi"
    "ndstore/dsrpc"
)

func (contr *Contr) AddBStoreHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewAddBStoreParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    address := params.Address
    port    := params.Port
    login   := params.Login
    pass    := params.Pass
    userName := string(context.AuthIdent())
    err = contr.store.AddBStore(userName, address, port, login, pass)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := fsapi.NewAddBStoreResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) UpdateBStoreHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewUpdateBStoreParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    address := params.Address
    port    := params.Port
    login   := params.Login
    pass    := params.Pass
    userName := string(context.AuthIdent())
    err = contr.store.UpdateBStore(userName, address, port, login, pass)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := fsapi.NewUpdateBStoreResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) DeleteBStoreHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewDeleteBStoreParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    address := params.Address
    port    := params.Port
    userName := string(context.AuthIdent())
    err = contr.store.DeleteBStore(userName, address, port)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fsapi.NewDeleteBStoreResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListBStoresHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewListBStoresParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    userName := string(context.AuthIdent())
    bstores, err := contr.store.ListBStores(userName)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fsapi.NewListBStoresResult()
    result.BStores = bstores
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
