/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsucont

import (
    //"errors"
    //"io"
    "ndstore/bstore/bsapi"
    "ndstore/dsrpc"
    //"ndstore/dslog"
)

func (contr *Contr) AddUserHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewAddUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    pass    := params.Pass
    role    := params.Role
    state   := params.State
    userName := string(context.AuthIdent())
    err = contr.auth.AddUser(userName, login, pass, role, state)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := bsapi.NewAddUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) CheckUserHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewCheckUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    pass    := params.Pass
    userName := string(context.AuthIdent())

    //err = context.ReadBin(io.Discard)
    //if err != nil {
    //    context.SendError(err)
    //    return err
    //}

    ok, err := contr.auth.CheckUser(userName, login, pass)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := bsapi.NewCheckUserResult()
    result.Match = ok
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) UpdateUserHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewUpdateUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    pass    := params.Pass
    role    := params.Role
    state   := params.State
    userName := string(context.AuthIdent())

    err = contr.auth.UpdateUser(userName, login, pass, role, state)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := bsapi.NewUpdateUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) DeleteUserHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewDeleteUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    userName := string(context.AuthIdent())

    err = contr.auth.DeleteUser(userName, login)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewDeleteUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListUsersHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewListUsersParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    userName    := string(context.AuthIdent())
    users, err := contr.auth.ListUsers(userName)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := bsapi.NewListUsersResult()
    result.Users = users
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
