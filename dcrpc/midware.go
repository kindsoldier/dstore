package dcrpc

import (
    "time"
)

func LogRequest(context *Context) error {
    var err error
    logDebug("request:", string(context.reqBody.JSON()))
    return err
}

func LogResponse(context *Context) error {
    var err error
    logDebug("response:", string(context.resBody.JSON()))
    return err
}

func LogAccess(context *Context) error {
    var err error
    execTime := time.Now().Sub(context.start)
    logAccess(context.remoteHost, context.reqBody.Method, execTime)
    return err
}
