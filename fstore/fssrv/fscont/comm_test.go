package fdcont

import (
    "testing"
    "github.com/stretchr/testify/assert"

    "ndstore/fstore/fsapi"
    "ndstore/fstore/fssrv/fsrec"
    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dsrpc"
)


func Test_File_Hello(t *testing.T) {
    var err error

    dbPath := "postgres://pgsql@localhost/test"
    reg := fsreg.NewReg()
    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)
    err = reg.MigrateDB()
    assert.NoError(t, err)

    store := fsrec.NewStore(t.TempDir(), reg)
    contr := NewContr(store)

    helloResp := GetHelloMsg
    params := fsapi.NewGetHelloParams()
    params.Message = GetHelloMsg
    result := fsapi.NewGetHelloResult()
    err = dsrpc.LocalExec(fsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)
    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
