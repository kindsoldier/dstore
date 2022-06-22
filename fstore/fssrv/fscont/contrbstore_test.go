/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "testing"
    "github.com/stretchr/testify/require"

    "ndstore/fstore/fsapi"
    "ndstore/fstore/fssrv/fsrec"
    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dsrpc"
)


func Test_BStore_AddCheckDelete(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    dbPath := "postgres://pgsql@localhost/test"
    reg := fsreg.NewReg()
    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    model := fsrec.NewStore(rootDir, reg)
    require.NoError(t, err)

    err = model.SeedUsers()
    require.NoError(t, err)

    contr := NewContr(model)

    addParams := fsapi.NewAddBStoreParams()
    addParams.Address  = "127.0.0.1"
    addParams.Port     = "1234"
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456"
    addResult := fsapi.NewAddBStoreResult()

    auth := dsrpc.CreateAuth([]byte("admin"), []byte("admin"))

    err = dsrpc.LocalExec(fsapi.AddBStoreMethod, addParams, addResult, auth, contr.AddBStoreHandler)
    require.NoError(t, err)

    addParams = fsapi.NewAddBStoreParams()
    addParams.Address  = "127.0.0.1"
    addParams.Port     = "1234"
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456xxx"
    addResult = fsapi.NewAddBStoreResult()
    err = dsrpc.LocalExec(fsapi.AddBStoreMethod, addParams, addResult, auth, contr.AddBStoreHandler)
    require.Error(t, err)

    deleteParams := fsapi.NewDeleteBStoreParams()
    deleteParams.Address = "127.0.0.1"
    deleteParams.Port    = "1234"
    deleteResult := fsapi.NewDeleteBStoreResult()
    err = dsrpc.LocalExec(fsapi.DeleteBStoreMethod, deleteParams, deleteResult, auth, contr.DeleteBStoreHandler)
    require.NoError(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}
