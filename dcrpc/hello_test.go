package dcrpc

import (
    "encoding/json"
    "errors"
    "testing"
    "time"
    "log"
)


func TestHello(t *testing.T) {
    go servHello()
    time.Sleep(1 * time.Second)
    execHello()
}

func BenchmarkHello(b *testing.B) {
    go servHello()
    time.Sleep(1 * time.Second)
    b.ResetTimer()

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            execHello()
        }
    }
    b.SetParallelism(20)
    b.RunParallel(pBench)
}


func execHello() error {
    var err error
    params := NewHelloParams()
    params.Message = "hello rdrpc!"
    result := NewHelloResult()
    auth := CreateAuth([]byte("qwert"), []byte("12345"))
    err = Exec("127.0.0.1:8081", HelloMethod, params, result, auth)
    if err != nil {
        log.Println("server err:", err)
        return err
    }
    log.Println("server result:", string(result.JSON()))
    return err
}


func servHello() error {
    var err error

    //SetAccessWriter(io.Discard)
    //SetMessageWriter(io.Discard)
    serv := NewService()

    cont := NewController()
    serv.Handler(HelloMethod, cont.HelloHandler)

    mw := NewMiddleware()
    serv.PreMiddleware(mw.Auth)

    serv.PreMiddleware(LogRequest)
    serv.PostMiddleware(LogResponse)
    serv.PostMiddleware(LogAccess)

    err = serv.Listen(":8081")
    if err != nil {
        return err
    }
    return err
}


type Middleware struct {
}

func NewMiddleware() *Middleware {
    return &Middleware{}
}

func (cont *Middleware) Auth(context *Context) error {
    var err error
    reqIdent := context.AuthIdent()
    reqSalt := context.AuthSalt()
    reqHash := context.AuthHash()

    ident := reqIdent
    pass := []byte("12345")

    auth := context.Auth()
    log.Println("auth ", string(auth.JSON()))

    ok := CheckHash(ident, pass, reqSalt, reqHash)
    log.Println("auth ok:", ok)
    if !ok {
        err = errors.New("auth ident or pass missmatch")
        context.SendError(err)
        return err
    }

    return err
}

type Controller struct {
}

func NewController() *Controller {
    return &Controller{}
}

func (cont *Controller) HelloHandler(context *Context) error {
    var err error
    params := NewHelloParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    //log.Println("hello params:", string(params.JSON()))

    result := NewHelloResult()
    result.Message = "hello!"
    err = context.SendResult(result)
    if err != nil {
        return err
    }
    return err
}


const HelloMethod string = "hello"

type HelloParams struct {
    Message string      `json:"message" json:"message"`
}

func NewHelloParams() *HelloParams {
    return &HelloParams{}
}

func (this *HelloParams) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}


type HelloResult struct {
    Message string      `json:"message" json:"message"`
}

func NewHelloResult() *HelloResult {
    return &HelloResult{}
}

func (this *HelloResult) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}
