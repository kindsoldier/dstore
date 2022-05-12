package dcrpc

import (
    "encoding/json"
)

// Packet is used to store intermediate data
// on Context for debugging purposes
type Packet struct {
    header      []byte
    rcpPayload  []byte
}

func NewPacket() *Packet {
    return &Packet{}
}

func (this *Packet) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}
