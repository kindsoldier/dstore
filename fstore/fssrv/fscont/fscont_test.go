package fscont

import (
    "bytes"
    "math/rand"
    "testing"

    "ndstore/fstore/fsapi"
    "ndstore/fstore/fssrv/fsrec"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)


func TestHello(t *testing.T) {
    var err error
    helloResp := "hello!"

    params := fsapi.NewHelloParams()
    params.Message = "hello server!"

    result := fsapi.NewHelloResult()

    contr := NewContr()
    store := fsrec.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    err = dsrpc.LocalExec(fsapi.HelloMethod, params, result, nil, contr.HelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}

func TestSLD(t *testing.T) {
    var err error

    params := fsapi.NewSaveParams()

    params.ClusterId    = 1
    params.FileId       = 2
    params.BatchId      = 3
    params.BlockId      = 4
    result := fsapi.NewSaveResult()

    contr := NewContr()
    store := fsrec.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    data := make([]byte, 1024 * 1024)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dsrpc.LocalPut(fsapi.SaveMethod, reader, size, params, result, nil, contr.SaveHandler)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = dsrpc.LocalGet(fsapi.LoadMethod, writer, params, result, nil, contr.LoadHandler)
    assert.NoError(t, err)
    assert.Equal(t, len(data), len(writer.Bytes()))
    assert.Equal(t, data, writer.Bytes())

    err = dsrpc.LocalExec(fsapi.DeleteMethod, params, result, nil, contr.DeleteHandler)
    assert.NoError(t, err)
}
