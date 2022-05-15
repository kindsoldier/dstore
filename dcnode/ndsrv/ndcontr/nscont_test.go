package ndcontr

import (
    "bytes"
    "testing"
    "dcstore/dcnode/ndapi"
    "dcstore/dcrpc"
    "dcstore/dcnode/ndsrv/ndstore"

    "github.com/stretchr/testify/assert"
)


func TestHello(t *testing.T) {
    var err error
    helloResp := "hello!"

    params := ndapi.NewHelloParams()
    params.Message = "hello server!"

    result := ndapi.NewHelloResult()


    contr := NewContr()
    store := ndstore.NewStore("/tmp")
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    err = dcrpc.LocalExec(ndapi.HelloMethod, params, result, nil, contr.HelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}

func TestSave(t *testing.T) {
    var err error

    params := ndapi.NewSaveParams()
    params.ClusterId    = 1
    params.FileId       = 2
    params.BatchId      = 3
    params.BlockId      = 4
    result := ndapi.NewSaveResult()

    contr := NewContr()
    store := ndstore.NewStore("/tmp")
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    data := []byte("qwerty")
    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dcrpc.LocalPut(ndapi.SaveMethod, reader, size, params, result, nil, contr.SaveHandler)
    assert.NoError(t, err)
}
