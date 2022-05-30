/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "errors"
    "io"
    "ndstore/fstore/fsapi"
    //"ndstore/fstore/fssrv/fsrec"
    "ndstore/dsrpc"
)


func (contr *Contr) AuthMidware(context *dsrpc.Context) error {
    var err error
    login := context.AuthIdent()
    salt := context.AuthSalt()
    hash := context.AuthHash()

    usersDescr, exists, err := contr.store.GetUser(string(login))
    if err != nil {
        context.SendError(err)
        return err
    }
    if !exists {
        err = errors.New("login not exists")
        context.SendError(err)
        return err
    }

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return err
    }
    //auth := context.Auth()
    //dslog.LogDebug("auth ", string(auth.JSON()))

    pass := []byte(usersDescr.Pass)
    ok := dsrpc.CheckHash(login, pass, salt, hash)
    //dslog.LogDebug("auth ok:", ok)
    if !ok {
        err = errors.New("auth login or pass missmatch")
        context.SendError(err)
        return err
    }
    return err
}


func (contr *Contr) AddUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewAddUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    pass    := params.Pass
    err = contr.store.AddUser(login, pass)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := fsapi.NewAddUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) CheckUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewCheckUserParams()
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

    ok, err := contr.store.CheckUser(login, pass)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := fsapi.NewCheckUserResult()
    result.Match = ok
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) UpdateUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewUpdateUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    pass    := params.Pass

    err = contr.store.UpdateUser(login, pass)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := fsapi.NewUpdateUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) DeleteUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewDeleteUserParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    login   := params.Login
    err = contr.store.DeleteUser(login)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fsapi.NewDeleteUserResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListUsersHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewListUsersParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    users, err := contr.store.ListUsers()
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fsapi.NewListUsersResult()
    result.Users = users
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
