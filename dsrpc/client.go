/*
 *
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 *
 */

package dsrpc

import (
    "encoding/json"
    "errors"
    "io"
    "net"
)


func Put(address string, method string, reader io.Reader, size int64, param, result any, auth *Auth) error {
    var err error

    conn, err := net.Dial("tcp", address)
    if err != nil {
        return err
    }
    defer conn.Close()

    return ConnPut(conn, method, reader, size, param, result, auth)
}


func ConnPut(conn net.Conn, method string, reader io.Reader, size int64, param, result any, auth *Auth) error {
    var err error
    context := CreateContext(conn)
    context.reqRPC.Method = method
    context.reqRPC.Params = param
    context.reqRPC.Auth = auth
    context.resRPC.Result = result

    context.binReader = reader
    context.binWriter = conn

    context.reqHeader.binSize = size

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

    err = context.UploadBin()
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

func Get(address string, method string, writer io.Writer, param, result any, auth *Auth) error {
    var err error

    conn, err := net.Dial("tcp", address)
    if err != nil {
        return err
    }
    defer conn.Close()

    return ConnGet(conn, method, writer, param, result, auth)
}

func ConnGet(conn net.Conn, method string, writer io.Writer, param, result any, auth *Auth) error {
    var err error

    context := CreateContext(conn)
    context.reqRPC.Method = method
    context.reqRPC.Params = param
    context.reqRPC.Auth = auth
    context.resRPC.Result = result

    context.binReader = conn
    context.binWriter = writer

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
    err = context.ReadResponse()
    if err != nil {
        return err
    }
    err = context.DownloadBin()
    if err != nil {
        return err
    }
    err = context.BindResponse()
    if err != nil {
        return err
    }
    return err
}

func Exec(address, method string, param any, result any, auth *Auth) error {
    var err error

    conn, err := net.Dial("tcp", address)
    if err != nil {
        return err
    }
    defer conn.Close()

    err = ConnExec(conn, method, param, result, auth)
    if err != nil {
        return err
    }
    return err
}


func ConnExec(conn net.Conn, method string, param any, result any, auth *Auth) error {
    var err error

    context := CreateContext(conn)
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

    context.reqPacket.rcpPayload, err = context.reqRPC.Pack()
    if err != nil {
        return err
    }
    rpcSize := int64(len(context.reqPacket.rcpPayload))
    context.reqHeader.rpcSize = rpcSize

    context.reqPacket.header, err = context.reqHeader.Pack()
    if err != nil {
        return err
    }
    return err
}

func (context *Context) WriteRequest() error {
    var err error
    _, err = context.sockWriter.Write(context.reqPacket.header)
    if err != nil {
        return err
    }
    _, err = context.sockWriter.Write(context.reqPacket.rcpPayload)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) UploadBin() error {
    var err error
    _, err = CopyBytes(context.binReader, context.binWriter, context.reqHeader.binSize)
    return err
}

func (context *Context) ReadResponse() error {
    var err error

    context.resPacket.header, err = ReadBytes(context.sockReader, headerSize)
    if err != nil {
        return err
    }
    context.resHeader, err = UnpackHeader(context.resPacket.header)
    if err != nil {
        return err
    }
    rpcSize := context.resHeader.rpcSize
    context.resPacket.rcpPayload, err = ReadBytes(context.sockReader, rpcSize)
    if err != nil {
        return err
    }
    return err
}

func (context *Context) DownloadBin() error {
    var err error
    _, err = CopyBytes(context.binReader, context.binWriter, context.resHeader.binSize)
    return err
}

func (context *Context) BindResponse() error {
    var err error

    err = json.Unmarshal(context.resPacket.rcpPayload, context.resRPC)
    if err != nil {
        return err
    }
    if len(context.resRPC.Error) > 0 {
        return errors.New(context.resRPC.Error)
    }
    return err
}