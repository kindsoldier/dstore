/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fscont

import (
    "dstore/fstore/fsapi"
    "dstore/dsrpc"
    "dstore/dserr"
    "dstore/dsdescr"
)

func (contr *Contr) AddUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewAddUserParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    descr := dsdescr.NewUser()
    descr.Login   = params.Login
    descr.Pass    = params.Pass
    authLogin := string(context.AuthIdent())
    err = contr.store.AddUser(authLogin, descr)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    result := fsapi.NewAddUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) CheckUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewCheckUserParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    login       := params.Login
    pass        := params.Pass
    authLogin    := string(context.AuthIdent())
    ok, err := contr.store.CheckUser(authLogin, login, pass)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewCheckUserResult()
    result.Match = ok
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) UpdateUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewUpdateUserParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    descr := dsdescr.NewUser()
    descr.Login   = params.Login
    descr.Pass    = params.Pass
    descr.State   = ""   // todo
    descr.Role    = ""   // todo
    authLogin    := string(context.AuthIdent())
    err = contr.store.UpdateUser(authLogin, descr)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    result := fsapi.NewUpdateUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) DeleteUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewDeleteUserParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    login   := params.Login
    authLogin    := string(context.AuthIdent())
    err = contr.store.DeleteUser(authLogin, login)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewDeleteUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) ListUsersHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewListUsersParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    authLogin := string(context.AuthIdent())
    users, err := contr.store.ListUsers(authLogin)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fsapi.NewListUsersResult()
    result.Users = users
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
