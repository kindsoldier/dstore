package ndcontr

import (
    "bytes"
    "math/rand"
    "testing"

    "dcstore/dcnode/ndapi"
    "dcstore/dcnode/ndsrv/ndstore"
    "dcstore/dcrpc"

    "github.com/stretchr/testify/assert"
)


func TestHello(t *testing.T) {
    var err error
    helloResp := "hello!"

    params := ndapi.NewHelloParams()
    params.Message = "hello server!"

    result := ndapi.NewHelloResult()


    contr := NewContr()
    store := ndstore.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    err = dcrpc.LocalExec(ndapi.HelloMethod, params, result, nil, contr.HelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}

func TestSLD(t *testing.T) {
    var err error

    params := ndapi.NewSaveParams()    assert.Equal(t, len(data), len(writer.Bytes()))

    params.ClusterId    = 1
    params.FileId       = 2
    params.BatchId      = 3
    params.BlockId      = 4
    result := ndapi.NewSaveResult()

    contr := NewContr()
    store := ndstore.NewStore(t.TempDir())
    err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    data := make([]byte, 1024 * 1024)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dcrpc.LocalPut(ndapi.SaveMethod, reader, size, params, result, nil, contr.SaveHandler)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = dcrpc.LocalGet(ndapi.LoadMethod, writer, params, result, nil, contr.LoadHandler)
    assert.NoError(t, err)
    assert.Equal(t, len(data), len(writer.Bytes()))
    assert.Equal(t, data, writer.Bytes())

    err = dcrpc.LocalExec(ndapi.DeleteMethod, params, result, nil, contr.DeleteHandler)
    assert.NoError(t, err)
}
