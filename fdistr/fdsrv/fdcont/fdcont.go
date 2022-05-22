/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "ndstore/fdistr/fdapi"
    "ndstore/dsrpc"
)

type Contr struct {
}

func NewContr() *Contr {
    return &Contr{}
}

func (contr *Contr) HelloHandler(context *dsrpc.Context) error {
    var err error
    params := fdapi.NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := fdapi.NewHelloResult()
    result.Message = "hello!"
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}
