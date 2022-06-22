/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "errors"
    "io"
    "ndstore/fstore/fsapi"
    "ndstore/dsrpc"
    "ndstore/dslog"
    "ndstore/dserr"
)

func (contr *Contr) AuthMidware(context *dsrpc.Context) error {
    var err error
    login := context.AuthIdent()
    salt := context.AuthSalt()
    hash := context.AuthHash()

    usersDescr, err := contr.store.GetUser(string(login))
    if err != nil {
        context.ReadBin(io.Discard)
        extErr := errors.New("auth missmatch")
        context.SendError(extErr)
        return dserr.Err(err)
    }

    auth := context.Auth()
    dslog.LogDebug("auth ", string(auth.JSON()))

    pass := []byte(usersDescr.Pass)
    ok := dsrpc.CheckHash(login, pass, salt, hash)
    dslog.LogDebug("auth ok:", ok)

    if !ok {
        context.ReadBin(io.Discard)
        extErr := errors.New("auth missmatch")
        context.SendError(extErr)
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) AddUserHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewAddUserParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    login   := params.Login
    pass    := params.Pass
    userName := string(context.AuthIdent())
    err = contr.store.AddUser(userName, login, pass)
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
        return dserr.Err(err)
    }
    login       := params.Login
    pass        := params.Pass
    userName    := string(context.AuthIdent())
    ok, err := contr.store.CheckUser(userName, login, pass)
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
        return dserr.Err(err)
    }
    login       := params.Login
    pass        := params.Pass
    state       := ""   // todo
    role        := ""   // todo
    userName    := string(context.AuthIdent())
    err = contr.store.UpdateUser(userName, login, pass, role, state)
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
        return dserr.Err(err)
    }
    login   := params.Login
    userName    := string(context.AuthIdent())
    err = contr.store.DeleteUser(userName, login)
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
        return dserr.Err(err)
    }
    userName := string(context.AuthIdent())
    users, err := contr.store.ListUsers(userName)
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
