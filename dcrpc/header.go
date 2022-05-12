package dcrpc

import (
    "encoding/binary"
    "encoding/json"
    "bytes"
)

const headerSize    int64   = 16
const sizeOfInt64   int     = 8


type Header struct {
    rpcSize int64          `json:"rpcSize"`
    binSize int64          `json:"binSize"`
}

func NewHeader() *Header {
    return &Header{}
}

func (this *Header) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}


func (this *Header) Pack() ([]byte, error) {
    var err error
    headerBytes := make([]byte, 0, headerSize)
    headerBuffer := bytes.NewBuffer(headerBytes)

    rpcSizeBytes := encoderI64(this.rpcSize)
    headerBuffer.Write(rpcSizeBytes)

    binSizeBytes := encoderI64(this.binSize)
    headerBuffer.Write(binSizeBytes)

    return headerBuffer.Bytes(), err
}

func UnpackHeader(headerBytes []byte) (*Header, error) {
    var err error
    header := NewHeader()
    headerReader := bytes.NewReader(headerBytes)

    rcpSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(rcpSizeBytes)
    header.rpcSize = decoderI64(rcpSizeBytes)

    binSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(binSizeBytes)
    header.binSize = decoderI64(binSizeBytes)

    return header, err
}

func encoderI64(i int64) []byte {
    buffer := make([]byte, sizeOfInt64)
    binary.BigEndian.PutUint64(buffer, uint64(i))
    return buffer
}

func decoderI64(b []byte) int64 {
    return int64(binary.BigEndian.Uint64(b))
}
