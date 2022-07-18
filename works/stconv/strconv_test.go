package dsstconv

import (
    "testing"
    "github.com/stretchr/testify/require"
)


type Descr struct {
    StringVal   string      `json:"stringVal"`

    Int64Val    int64       `json:"int64Val"`
    Int32Val    int32       `json:"int32Val"`
    Int16Val    int16       `json:"int16Val"`
    Int8Val     int8        `json:"int8Val"`

    Uint64Val   uint64      `json:"uint64Val"`
    Uint32Val   uint32      `json:"uint32Val"`
    Uint16Val   uint16      `json:"uint16Val"`
    Uint8Val    uint8       `json:"uint8Val"`

    Float64Val  float64     `json:"float64Val"`
    Float32Val  float32     `json:"float32Val"`

    BoolVal     bool       `json:"boolVal"`
}

func NewDescr() *Descr {
    var descr Descr
    return &descr
}


func TestConv(t *testing.T)  {
    var err error

    descr1 := NewDescr()
    descr1.StringVal    = "qwerty"

    descr1.Int64Val     = -1
    descr1.Int32Val     = -2
    descr1.Int16Val     = -3
    descr1.Int8Val      = -4

    descr1.Uint64Val    = 1
    descr1.Uint32Val    = 2
    descr1.Uint16Val    = 3
    descr1.Uint8Val     = 4

    descr1.Float64Val   = 1.123456
    descr1.Float32Val   = 6.789012

    descr1.BoolVal      = true

    resMap, err := Pack(descr1)
    require.NoError(t, err)

    descr2 := NewDescr()
    err = Unpack(resMap, descr2)
    require.NoError(t, err)
    require.Equal(t, descr1, descr2)

}
