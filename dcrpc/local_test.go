package dcrpc

import (
    "bytes"
    "fmt"
    "testing"
)

func localHello() error {
    var err error

    reqBuffer := make([]byte, 0)
    reqRW := bytes.NewBuffer(reqBuffer)

    resBuffer := make([]byte, 0)
    resRW := bytes.NewBuffer(resBuffer)

    context := CreateLocalContext(resRW, reqRW)

    params := NewHelloParams()
    params.Message = "hello rdrpc!"
    result := NewHelloResult()
    auth := CreateAuth([]byte("qwert"), []byte("12345"))
    context.WriteLocalRequest(HelloMethod, params, result, auth)

    serv := NewService()

    cont := NewController()
    serv.Handler(HelloMethod, cont.HelloHandler)

    err = serv.HandleLocal(resRW, reqRW)
    if err != nil {
        return err
    }

    err = context.BindLocalResult()
    if err != nil {
        return err
    }
    fmt.Println("result:", string(result.JSON()))

    return err
}

func TestLocal(t *testing.T) {
    err := localHello()
    if err != nil {
        t.Fatal(err)
    }
}
