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


func Test_User_AddCheckDelete(t *testing.T) {
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

    auth := dsrpc.CreateAuth([]byte("admin"), []byte("admin"))

    addParams := fsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456"
    addResult := fsapi.NewAddUserResult()
    err = dsrpc.LocalExec(fsapi.AddUserMethod, addParams, addResult, auth, contr.AddUserHandler)
    require.NoError(t, err)

    checkParams := fsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult := fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.NoError(t, err)
    require.Equal(t, true, checkResult.Match)


    addParams = fsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456xx"
    addResult = fsapi.NewAddUserResult()
    err = dsrpc.LocalExec(fsapi.AddUserMethod, addParams, addResult, auth, contr.AddUserHandler)
    require.Error(t, err)

    addParams = fsapi.NewAddUserParams()
    addParams.Login    = "йцукен"
    addParams.Pass     = "567890"
    addResult = fsapi.NewAddUserResult()
    err = dsrpc.LocalExec(fsapi.AddUserMethod, addParams, addResult, auth, contr.AddUserHandler)
    require.NoError(t, err)


    checkParams = fsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456XXX"
    checkResult = fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.NoError(t, err)
    require.Equal(t, false, checkResult.Match)

    checkParams = fsapi.NewCheckUserParams()
    checkParams.Login    = "qwertyXXX"
    checkParams.Pass     = "123456"
    checkResult = fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.Error(t, err)
    require.Equal(t, false, checkResult.Match)

    deleteParams := fsapi.NewDeleteUserParams()
    deleteParams.Login    = "qwerty"
    deleteResult := fsapi.NewDeleteUserResult()
    err = dsrpc.LocalExec(fsapi.DeleteUserMethod, deleteParams, deleteResult, auth, contr.DeleteUserHandler)
    require.NoError(t, err)

    checkParams = fsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult = fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.Error(t, err)
    require.Equal(t, false, checkResult.Match)

    err = reg.CloseDB()
    require.NoError(t, err)

}
