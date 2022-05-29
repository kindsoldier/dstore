package bscont

import (
    "testing"
    "path/filepath"

    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsblock"
    "ndstore/bstore/bssrv/bsbreg"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)


func TestGetHello(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "blocks.db")
    reg := bsbreg.NewReg()
    err = reg.OpenDB(path)
    assert.NoError(t, err)
    err = reg.MigrateDB()
    assert.NoError(t, err)

    store := bsblock.NewStore(rootDir, reg)
    assert.NoError(t, err)

    contr := NewContr(store)

    helloResp := GetHelloMsg
    params := bsapi.NewGetHelloParams()
    params.Message = GetHelloMsg
    result := bsapi.NewGetHelloResult()
    err = dsrpc.LocalExec(bsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}
