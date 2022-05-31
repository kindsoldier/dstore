/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsucont

import (
    "errors"
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/dsrpc"
    "ndstore/dslog"
)

func (contr *Contr) AuthMidware(context *dsrpc.Context) error {
    var err error
    login := context.AuthIdent()
    salt := context.AuthSalt()
    hash := context.AuthHash()

    usersDescr, err := contr.auth.GetUser(string(login))
    if err != nil {
        context.SendError(err)
        return err
    }
    auth := context.Auth()
    dslog.LogDebug("auth ", string(auth.JSON()))

    pass := []byte(usersDescr.Pass)
    ok := dsrpc.CheckHash(login, pass, salt, hash)
    dslog.LogDebug("auth ok:", ok)
    if !ok {
        context.ReadBin(io.Discard)

        err = errors.New("auth login or pass missmatch")
        context.SendError(err)
        return err
    }
    return err
}


func (contr *Contr) AddUserHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewAddUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    pass    := params.Pass
    err = contr.auth.AddUser(login, pass)
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

    //err = context.ReadBin(io.Discard)
    //if err != nil {
    //    context.SendError(err)
    //    return err
    //}

    ok, err := contr.auth.CheckUser(login, pass)
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

    err = contr.auth.UpdateUser(login, pass)
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
    err = contr.auth.DeleteUser(login)
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
    users, err := contr.auth.ListUsers()
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
