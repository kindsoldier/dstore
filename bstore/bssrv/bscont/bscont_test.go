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


func TestGetHello(t *testing.T) {
    var err error
    helloResp := GetHelloMsg

    params := bsapi.NewGetHelloParams()
    params.Message = GetHelloMsg

    result := bsapi.NewGetHelloResult()

    contr := NewContr()
    store := bsrec.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    err = dsrpc.LocalExec(bsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}

func TestSaveLoadDelete(t *testing.T) {
    var err error

    params := bsapi.NewSaveBlockParams()

    params.FileId       = 2
    params.BatchId      = 3
    params.BlockId      = 4
    result := bsapi.NewSaveBlockResult()

    contr := NewContr()
    store := bsrec.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    data := make([]byte, 1024 * 1024)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dsrpc.LocalPut(bsapi.SaveBlockMethod, reader, size, params, result, nil, contr.SaveBlockHandler)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = dsrpc.LocalGet(bsapi.LoadBlockMethod, writer, params, result, nil, contr.LoadBlockHandler)
    assert.NoError(t, err)
    assert.Equal(t, len(data), len(writer.Bytes()))
    assert.Equal(t, data, writer.Bytes())

    err = dsrpc.LocalExec(bsapi.DeleteBlockMethod, params, result, nil, contr.DeleteBlockHandler)
    assert.NoError(t, err)
}
