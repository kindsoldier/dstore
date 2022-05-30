/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "ndstore/fstore/fsapi"
    "ndstore/fstore/fssrv/fsrec"
    "ndstore/dsrpc"
)


type Contr struct {
    store  *fsrec.Store
}

func NewContr(store *fsrec.Store) *Contr {
    return &Contr{ store: store }
}

const GetHelloMsg string = "hello"

func (contr *Contr) GetHelloHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewGetHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    result := fsapi.NewGetHelloResult()
    result.Message = GetHelloMsg
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
