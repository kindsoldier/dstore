package dcrpc

import (
    "io"
    "net"
    "time"
)

type Context struct {
    start       time.Time
    remoteHost  string
    reader      io.Reader
    writer      io.Writer

    reqPacket   *Packet
    reqHeader   *Header
    reqBody     *Request

    resPacket   *Packet
    resHeader   *Header
    resBody     *Response
}


func NewContext() *Context {
    context := &Context{}
    context.start = time.Now()
    return context
}


func CreateContext(conn net.Conn) *Context {
    context := &Context{}
    context.start = time.Now()
    context.reader = conn
    context.writer = conn

    context.reqPacket = NewPacket()
    context.resPacket = NewPacket()

    context.reqHeader = NewHeader()
    context.reqBody   = NewRequest()

    context.resHeader = NewHeader()
    context.resBody = NewResponse()
    context.resBody = NewResponse()

    return context
}

func (context *Context) Request() *Request  {
    return context.reqBody
}


func (context *Context) SetAuthIdent(ident []byte)  {
    context.reqBody.Auth.Ident = ident
}

func (context *Context) SetAuthSalt(salt []byte)  {
    context.reqBody.Auth.Salt = salt
}

func (context *Context) SetAuthHash(hash []byte)  {
    context.reqBody.Auth.Hash = hash
}

func (context *Context) AuthIdent() []byte {
    return context.reqBody.Auth.Ident
}

func (context *Context) AuthSalt() []byte {
    return context.reqBody.Auth.Salt
}

func (context *Context) AuthHash() []byte {
    return context.reqBody.Auth.Hash
}

func (context *Context) Auth() *Auth {
    return context.reqBody.Auth
}
