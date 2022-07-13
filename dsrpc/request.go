/*
 *
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 *
 */

package dsrpc

import (
    "encoding/json"
    "github.com/shamaton/msgpack/v2"
)

type Request struct {
    Method  string      `json:"method"  msgpack:"method"`
    Params  any         `json:"params"  msgpack:"params"`
    Auth    *Auth       `json:"auth"    msgpack:"auth"`
}

func NewRequest() *Request {
    req := &Request{}
    req.Auth = &Auth{}
    return req
}

func (this *Request) Pack() ([]byte, error) {
    rBytes, err := msgpack.Marshal(this)
    return rBytes, Err(err)
}

func (this *Request) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}
