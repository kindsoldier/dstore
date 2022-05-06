package dcrpc

import (
    "encoding/binary"
    "encoding/json"
    "bytes"
)

const headerSize    int64   = 8
const sizeOfInt64   int     = 8


type Header struct {
    BodySize int64          `json:"bodySize"`

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

    bodySizeBytes := encoderI64(this.BodySize)
    headerBuffer.Write(bodySizeBytes)

    return headerBuffer.Bytes(), err
}

func UnpackHeader(headerBytes []byte) (*Header, error) {
    var err error
    header := NewHeader()
    headerReader := bytes.NewReader(headerBytes)

    bodySizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(bodySizeBytes)
    header.BodySize = decoderI64(bodySizeBytes)

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
