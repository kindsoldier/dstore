package bscont

import (
    "bytes"
    "math/rand"
    "testing"

    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsrec"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)


func TestHello(t *testing.T) {
    var err error
    helloResp := HelloMsg

    params := bsapi.NewHelloParams()
    params.Message = HelloMsg

    result := bsapi.NewHelloResult()

    contr := NewContr()
    store := bsrec.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    err = dsrpc.LocalExec(bsapi.HelloMethod, params, result, nil, contr.HelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}

func TestSLD(t *testing.T) {
    var err error

    params := bsapi.NewSaveParams()

    params.ClusterId    = 1
    params.FileId       = 2
    params.BatchId      = 3
    params.BlockId      = 4
    result := bsapi.NewSaveResult()

    contr := NewContr()
    store := bsrec.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    data := make([]byte, 1024 * 1024)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dsrpc.LocalPut(bsapi.SaveMethod, reader, size, params, result, nil, contr.SaveHandler)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = dsrpc.LocalGet(bsapi.LoadMethod, writer, params, result, nil, contr.LoadHandler)
    assert.NoError(t, err)
    assert.Equal(t, len(data), len(writer.Bytes()))
    assert.Equal(t, data, writer.Bytes())

    err = dsrpc.LocalExec(bsapi.DeleteMethod, params, result, nil, contr.DeleteHandler)
    assert.NoError(t, err)
}
