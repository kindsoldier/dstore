/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "dcstore/dcnode/nodeapi"
    "dcstore/dclog"
    "dcstore/rdrpc"
)


type Controller struct {
}

func NewController() *Controller {
    return &Controller{}
}

func (cont *Controller) HelloHandler(context *rdrpc.Context) error {
    var err error
    params := ndapi.NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    dclog.LogDebug("hello params:", string(params.JSON()))

    result := ndapi.NewHelloResult()
    result.Message = "hello!"
    err = context.SendResult(result)
    if err != nil {
        return err
    }
    return err
}
