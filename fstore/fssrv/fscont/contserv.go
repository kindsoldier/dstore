/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fscont

import (
    "dstore/fstore/fsapi"
    "dstore/dsrpc"
    "dstore/dserr"
)

func (contr *Contr) GetStatusHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewGetStatusParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    result := fsapi.NewGetStatusResult()
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
