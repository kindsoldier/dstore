package fdcont

import (
    "testing"

    "ndstore/fstore/fsapi"
    "ndstore/fstore/fssrv/fsrec"
    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)


func Test_File_Hello(t *testing.T) {
    var err error

    dbPath := "postgres://pgsql@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    helloResp := GetHelloMsg

    params := fsapi.NewGetHelloParams()
    params.Message = GetHelloMsg

    result := fsapi.NewGetHelloResult()

    contr := NewContr()
    store := fsrec.NewStore(t.TempDir(), reg)
    assert.NoError(t, err)

    contr.Store = store

    err = dsrpc.LocalExec(fsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
