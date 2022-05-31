/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsucont

import (
    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsuser"
    "ndstore/dsrpc"
)

const GetHelloMsg string = "hello"

type Contr struct {
    auth    *bsuser.Auth
}

func NewContr(auth *bsuser.Auth) *Contr {
    return &Contr{ auth: auth }
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
