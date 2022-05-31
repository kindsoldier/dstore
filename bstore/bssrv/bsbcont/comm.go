/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsbcont

import (
    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsblock"
    "ndstore/dsrpc"
)

type Contr struct {
    store   *bsblock.Store
}

func NewContr(store *bsblock.Store) *Contr {
    return &Contr{ store: store }
}

const GetHelloMsg string = "hello"

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
