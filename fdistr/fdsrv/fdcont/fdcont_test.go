package fdcont

import (
    "bytes"
    "math/rand"
    "testing"

    "ndstore/fdistr/fdapi"
    "ndstore/fdistr/fdsrv/fdrec"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)


func TestHello(t *testing.T) {
    var err error
    helloResp := HelloMsg

    params := fdapi.NewHelloParams()
    params.Message = HelloMsg

    result := fdapi.NewHelloResult()

    contr := NewContr()
    store := fdrec.NewStore(t.TempDir())
    //err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    err = dsrpc.LocalExec(fdapi.HelloMethod, params, result, nil, contr.HelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}

func TestSaveLoadDelete(t *testing.T) {
    var err error

    params := fdapi.NewSaveParams()
    params.FilePath = "qwert.txt"

    result := fdapi.NewSaveResult()

    contr := NewContr()
    store := fdrec.NewStore(t.TempDir())
    //err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    data := make([]byte, 1024 * 1024)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dsrpc.LocalPut(fdapi.SaveMethod, reader, size, params, result, nil, contr.SaveHandler)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = dsrpc.LocalGet(fdapi.LoadMethod, writer, params, result, nil, contr.LoadHandler)
    assert.NoError(t, err)
    assert.Equal(t, len(data), len(writer.Bytes()))
    assert.Equal(t, data, writer.Bytes())

    err = dsrpc.LocalExec(fdapi.DeleteMethod, params, result, nil, contr.DeleteHandler)
    assert.NoError(t, err)
}
