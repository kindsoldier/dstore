package dcrpc

import (
    "context"
    "errors"
    "net"
    "sync"

    "encoding/json"
)

type HandlerFunc =  func(*Context) error

type Service struct {
    handlers map[string]HandlerFunc
    ctx     context.Context
    cancel  context.CancelFunc
    wg      *sync.WaitGroup
    preMw   []HandlerFunc
    postMw  []HandlerFunc
}

func NewService() *Service {
    rdrpc := &Service{}
    rdrpc.handlers = make(map[string]HandlerFunc)
    ctx, cancel := context.WithCancel(context.Background())
    rdrpc.ctx = ctx
    rdrpc.cancel = cancel
    var wg sync.WaitGroup
    rdrpc.wg = &wg
    rdrpc.preMw = make([]HandlerFunc, 0)
    rdrpc.postMw = make([]HandlerFunc, 0)

    return rdrpc
}

func (this *Service) PreMiddleware(mw HandlerFunc) {
    this.preMw = append(this.preMw, mw)
}

func (this *Service) PostMiddleware(mw HandlerFunc) {
    this.postMw = append(this.postMw, mw)
}


func (this *Service) Handler(method string, handler HandlerFunc) {
    this.handlers[method] = handler
}

func (this *Service) Listen(address string) error {
    var err error
    logInfo("server listen:", address)
    listener, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }
    this.wg.Add(1)
    for {
        select {
            case <- this.ctx.Done():
                this.wg.Done()
                return err
            default:
        }
        conn, err := listener.Accept()
        if err != nil {
            logError("conn accept err:", err)
        }

        go this.handleConn(conn)
    }
}

func notFound(context *Context) error {
    execErr := errors.New("method not found")
    err := context.SendError(execErr)
    return err
}

func (this *Service) Stop() error {
    var err error
    this.cancel()
    this.wg.Wait()
    return err
}

func (this *Service) handleConn(conn net.Conn) {
    var err error

    context := CreateContext(conn)

    remoteAddr := conn.RemoteAddr().String()
    remoteHost, _, _ := net.SplitHostPort(remoteAddr)
    context.remoteHost = remoteHost

    exitFunc := func() {
            conn.Close()
            if err != nil {
                logError("conn handler err:", err)
            }
    }
    defer exitFunc()

    err = context.ReadRequest()
    if err != nil {
        return
    }
    err = context.BindMethod()
    if err != nil {
        return
    }
    for _, mw := range this.preMw {
        err = mw(context)
        if err != nil {
            return
        }
    }
    err = this.Route(context)
    if err != nil {
        return
    }
    for _, mw := range this.postMw {
        err = mw(context)
        if err != nil {
            return
        }
    }
    return
}

func (this *Service) Route(context *Context) error {
    handler, ok := this.handlers[context.reqBody.Method]
    if ok {
        return handler(context)
    }
    return notFound(context)
}

func (context *Context) ReadRequest() error {
    var err error

    context.reqPacket.header, err = ReadBytes(context.reader, headerSize)
    if err != nil {
        return err
    }
    context.reqHeader, err = UnpackHeader(context.reqPacket.header)
    if err != nil {
        return err
    }
    bodySize := context.reqHeader.BodySize
    context.reqPacket.body, err = ReadBytes(context.reader, bodySize)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) BindMethod() error {
    var err error
    err = json.Unmarshal(context.reqPacket.body, context.reqBody)
    return err
}

func (context *Context) BindParams(params any) error {
    var err error
    context.reqBody.Params = params
    err = json.Unmarshal(context.reqPacket.body, context.reqBody)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) SendResult(result any) error {
    var err error
    context.resBody.Result = result

    context.resPacket.body, err = context.resBody.Pack()
    if err != nil {
        return err
    }
    context.resHeader.BodySize = int64(len(context.resPacket.body))

    context.resPacket.header, err = context.resHeader.Pack()
    if err != nil {
        return err
    }
    _, err = context.writer.Write(context.resPacket.header)
    if err != nil {
        return err
    }
    _, err = context.writer.Write(context.resPacket.body)
    if err != nil {
        return err
    }
    return err
}


func (context *Context) SendError(execErr error) error {
    var err error

    context.resBody.Error = execErr.Error()
    context.resBody.Result = NewEmpty()

    context.resPacket.body, err = context.resBody.Pack()
    if err != nil {
        return err
    }
    context.resHeader.BodySize = int64(len(context.resPacket.body))
    context.resPacket.header, err = context.resHeader.Pack()
    if err != nil {
        return err
    }
    _, err = context.writer.Write(context.resPacket.header)
    if err != nil {
        return err
    }
    _, err = context.writer.Write(context.resPacket.body)
    if err != nil {
        return err
    }
    return err
}
