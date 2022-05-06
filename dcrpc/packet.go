package dcrpc

import (
    "encoding/json"
)

type Packet struct {
    header  []byte
    body    []byte
}

func NewPacket() *Packet {
    return &Packet{}
}

func (this *Packet) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}
