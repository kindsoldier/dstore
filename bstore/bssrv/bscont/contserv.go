/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bscont

import (
    "dstore/bstore/bsapi"
    "dstore/dscomm/dsrpc"
    "dstore/dscomm/dserr"
)

func (contr *Contr) GetStatusHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewGetStatusParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewGetStatusResult()
    uptime, err := contr.store.GetUptime()
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result.Uptime = uptime
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
