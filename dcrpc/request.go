package dcrpc

import (
    "encoding/json"
)

type Request struct {
    Method  string      `json:"method"           msgpack:"method"`
    Params  any         `json:"params,omitempty" msgpack:"params,omitempty"`
    Auth    *Auth       `json:"auth,omitempty"   msgpack:"auth,omitempty"`
}

func NewRequest() *Request {
    req := &Request{}
    req.Auth = &Auth{}
    return req
}

func (this *Request) Pack() ([]byte, error) {
    rBytes, err := json.Marshal(this)
    return rBytes, err
}

func (this *Request) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}
