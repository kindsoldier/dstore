/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsucont

import (
    "errors"
    "io"
    "ndstore/bstore/bssrv/bsuser"
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

    if usersDescr.State != bsuser.UStateEnabled {
        context.ReadBin(io.Discard)

        err = errors.New("user state is not enabled")
        context.SendError(err)
        return err
    }

    return err
}
