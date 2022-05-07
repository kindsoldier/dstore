package dcrpc

import (
    "io"
    "time"
)

func CreateLocalContext(writer io.Writer, reader io.Reader) *Context {
    context := &Context{}
    context.start = time.Now()

    context.writer = writer
    context.reader = reader

    context.reqPacket = NewPacket()
    context.resPacket = NewPacket()

    context.reqHeader = NewHeader()
    context.reqBody = NewRequest()

    context.resHeader = NewHeader()
    context.resBody = NewResponse()
    context.resBody = NewResponse()

    return context
}

func (context *Context) WriteLocalRequest(method string, param, result interface{}, auth *Auth) error {
    var err error

    context.reqBody.Method = method
    context.reqBody.Params = param
    context.reqBody.Auth = auth
    context.resBody.Result = result

    if context.reqBody.Params == nil {
        context.reqBody.Params = NewEmpty()
    }

    err = context.CreateRequest()
    if err != nil {
        return err
    }
    err = context.WriteRequest()
    if err != nil {
        return err
    }
    return err
}

func (this *Service) HandleLocal(reqReader io.Reader, resWriter io.Writer) error {
    var err error

    context := CreateLocalContext(resWriter, reqReader)
    const fakeRemoteHost = "101.111.111.111"
    context.remoteHost = fakeRemoteHost

    err = context.ReadRequest()
    if err != nil {
        return err
    }
    err = context.BindMethod()
    if err != nil {
        return err
    }
    for _, mw := range this.preMw {
        err = mw(context)
        if err != nil {
            return err
        }
    }
    err = this.Route(context)
    if err != nil {
        return err
    }
    for _, mw := range this.postMw {
        err = mw(context)
        if err != nil {
            return err
        }
    }
    return err
}


func (context *Context) BindLocalResult() error {
    var err error
    err = context.ReadResponse()
    if err != nil {
        return err
    }
    err = context.BindResponse()
    if err != nil {
        return err
    }
    return err
}
