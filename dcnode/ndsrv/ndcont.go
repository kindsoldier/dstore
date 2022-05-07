/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "dcstore/dcnode/ndapi"
    "dcstore/dcrpc"
)


type Controller struct {
}

func NewController() *Controller {
    return &Controller{}
}

func (cont *Controller) HelloHandler(context *dcrpc.Context) error {
    var err error
    params := ndapi.NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    result := ndapi.NewHelloResult()
    result.Message = "hello!"
    err = context.SendResult(result)
    if err != nil {
        return err
    }
    return err
}
