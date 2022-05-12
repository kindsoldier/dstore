package dcrpc

import (
    "io"
    "time"
)

func CreateLocalContext(writer io.Writer, reader io.Reader) *Context {
    context := &Context{}
    context.start = time.Now()

    context.sockWriter = writer
    context.sockReader = reader

    context.reqPacket = NewPacket()
    context.resPacket = NewPacket()

    context.reqHeader = NewHeader()
    context.reqRPC = NewRequest()

    context.resHeader = NewHeader()
    context.resRPC = NewResponse()
    context.resRPC = NewResponse()

    return context
}

func (context *Context) WriteLocalRequest(method string, param, result any, auth *Auth) error {
    var err error

    context.reqRPC.Method = method
    context.reqRPC.Params = param
    context.reqRPC.Auth = auth
    context.resRPC.Result = result

    if context.reqRPC.Params == nil {
        context.reqRPC.Params = NewEmpty()
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
