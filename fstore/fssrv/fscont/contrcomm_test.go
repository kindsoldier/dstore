package fdcont

import (
    "testing"
    "github.com/stretchr/testify/require"

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
    require.NoError(t, err)
    err = reg.MigrateDB()
    require.NoError(t, err)

    store := fsrec.NewStore(t.TempDir(), reg)
    contr := NewContr(store)

    helloResp := GetHelloMsg
    params := fsapi.NewGetHelloParams()
    params.Message = GetHelloMsg
    result := fsapi.NewGetHelloResult()
    err = dsrpc.LocalExec(fsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)
    require.NoError(t, err)
    require.Equal(t, helloResp, result.Message)

    err = reg.CloseDB()
    require.NoError(t, err)
}
