package fdcont

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
    helloResp := HelloMsg

    params := fsapi.NewHelloParams()
    params.Message = HelloMsg

    result := fsapi.NewHelloResult()

    contr := NewContr()
    store := fdrec.NewStore(t.TempDir())
    //err = store.OpenReg()
    assert.NoError(t, err)

    contr.Store = store

    err = dsrpc.LocalExec(fsapi.HelloMethod, params, result, nil, contr.HelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}

func TestSaveLoadDelete(t *testing.T) {
    var err error

    params := fsapi.NewSaveParams()
    params.FilePath = "qwert.txt"

    result := fsapi.NewSaveResult()

    contr := NewContr()
    store := fdrec.NewStore(t.TempDir())
    //err = store.OpenReg()
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
