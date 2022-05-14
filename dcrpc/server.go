package dcrpc

import (
    "context"
    "encoding/json"
    "errors"
    "io"
    "net"
    "sync"
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

    context.binReader = conn
    context.binWriter = io.Discard

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
    handler, ok := this.handlers[context.reqRPC.Method]
    if ok {
        return handler(context)
    }
    return notFound(context)
}

func (context *Context) ReadRequest() error {
    var err error

    context.reqPacket.header, err = ReadBytes(context.sockReader, headerSize)
    if err != nil {
        return err
    }
    context.reqHeader, err = UnpackHeader(context.reqPacket.header)
    if err != nil {
        return err
    }

    rpcSize := context.reqHeader.rpcSize
    context.reqPacket.rcpPayload, err = ReadBytes(context.sockReader, rpcSize)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) BinReader() io.Reader {
    return context.sockReader
}

func (context *Context) BinSize() int64 {
    return context.reqHeader.binSize
}

func (context *Context) ReadBin(writer io.Writer) error {
    var err error
    _, err = CopyBytes(context.sockReader, writer, context.reqHeader.binSize)
    return err
}


func (context *Context) BindMethod() error {
    var err error
    err = json.Unmarshal(context.reqPacket.rcpPayload, context.reqRPC)
    return err
}

func (context *Context) BindParams(params any) error {
    var err error
    context.reqRPC.Params = params
    err = json.Unmarshal(context.reqPacket.rcpPayload, context.reqRPC)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) SendResult(result any) error {
    var err error
    context.resRPC.Result = result

    context.resPacket.rcpPayload, err = context.resRPC.Pack()
    if err != nil {
        return err
    }
    context.resHeader.rpcSize = int64(len(context.resPacket.rcpPayload))

    context.resPacket.header, err = context.resHeader.Pack()
    if err != nil {
        return err
    }
    _, err = context.sockWriter.Write(context.resPacket.header)
    if err != nil {
        return err
    }
    _, err = context.sockWriter.Write(context.resPacket.rcpPayload)
    if err != nil {
        return err
    }
    return err
}


func (context *Context) SendBin(reader io.Reader, binSize int64, result any) error {
    var err error
    context.resRPC.Result = result

    context.resPacket.rcpPayload, err = context.resRPC.Pack()
    if err != nil {
        return err
    }
    context.resHeader.rpcSize = int64(len(context.resPacket.rcpPayload))
    context.resHeader.binSize = binSize

    context.resPacket.header, err = context.resHeader.Pack()
    if err != nil {
        return err
    }
    _, err = context.sockWriter.Write(context.resPacket.header)
    if err != nil {
        return err
    }
    _, err = context.sockWriter.Write(context.resPacket.rcpPayload)
    if err != nil {
        return err
    }

    _, err = CopyBytes(reader, context.sockWriter, context.resHeader.binSize)
    if err != nil {
        return err
    }

    return err
}


func (context *Context) SendError(execErr error) error {
    var err error

    context.resRPC.Error = execErr.Error()
    context.resRPC.Result = NewEmpty()

    context.resPacket.rcpPayload, err = context.resRPC.Pack()
    if err != nil {
        return err
    }
    context.resHeader.rpcSize = int64(len(context.resPacket.rcpPayload))
    context.resPacket.header, err = context.resHeader.Pack()
    if err != nil {
        return err
    }
    _, err = context.sockWriter.Write(context.resPacket.header)
    if err != nil {
        return err
    }
    _, err = context.sockWriter.Write(context.resPacket.rcpPayload)
    if err != nil {
        return err
    }
    return err
}
