package dcrpc

import (
    "errors"
    "net"
    "encoding/json"
)

func Exec(address, method string, param, result interface{}, auth *Auth) error {
    var err error

    conn, err := net.Dial("tcp", address)
    if err != nil {
        return err
    }
    defer conn.Close()

    context := CreateContext(conn)
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

func (context *Context) CreateRequest() error {
    var err error

    context.reqPacket.body, err = context.reqBody.Pack()
    if err != nil {
        return err
    }
    bodySize := int64(len(context.reqPacket.body))
    context.reqHeader.BodySize = bodySize

    context.reqPacket.header, err = context.reqHeader.Pack()
    if err != nil {
        return err
    }
    return err
}

func (context *Context) WriteRequest() error {
    var err error
    _, err = context.writer.Write(context.reqPacket.header)
    if err != nil {
        return err
    }
    _, err = context.writer.Write(context.reqPacket.body)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) ReadResponse() error {
    var err error

    context.resPacket.header, err = ReadBytes(context.reader, headerSize)
    if err != nil {
        return err
    }
    context.resHeader, err = UnpackHeader(context.resPacket.header)
    if err != nil {
        return err
    }
    bodySize := context.resHeader.BodySize
    context.resPacket.body, err = ReadBytes(context.reader, bodySize)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) BindResponse() error {
    var err error

    err = json.Unmarshal(context.resPacket.body, context.resBody)
    if err != nil {
        return err
    }
    if len(context.resBody.Error) > 0 {
        return errors.New(context.resBody.Error)
    }
    return err
}
