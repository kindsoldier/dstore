/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bscont

import (
    "errors"
    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsblock"
    "ndstore/bstore/bssrv/bsuser"
    "ndstore/dsrpc"
    "ndstore/dslog"
)

const GetHelloMsg string = "hello"

type Contr struct {
    store   *bsblock.Store
    auth    *bsuser.Auth
}

func NewContr(store *bsblock.Store, auth *bsuser.Auth) *Contr {
    return &Contr{ store: store, auth: auth }
}

func (contr *Contr) GetHelloHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewGetHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := bsapi.NewGetHelloResult()
    result.Message = GetHelloMsg
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) AuthMidware(context *dsrpc.Context) error {
    var err error
    login := context.AuthIdent()
    salt := context.AuthSalt()
    hash := context.AuthHash()

    usersDescr, exists, err := contr.auth.GetUser(string(login))
    if err != nil {
        context.SendError(err)
        return err
    }
    if !exists {
        err = errors.New("ident not exists")
        context.SendError(err)
        return err
    }

    auth := context.Auth()
    dslog.LogDebug("auth ", string(auth.JSON()))

    pass := []byte(usersDescr.Pass)
    ok := dsrpc.CheckHash(login, pass, salt, hash)
    dslog.LogDebug("auth ok:", ok)
    if !ok {
        err = errors.New("auth ident or pass missmatch")
        context.SendError(err)
        return err
    }
    return err
}
