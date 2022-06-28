/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsucont

import (
    "testing"
    "path/filepath"

    "github.com/stretchr/testify/require"

    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsuser"
    "ndstore/bstore/bssrv/bsureg"
    "ndstore/dsrpc"
)


func Test_User_AddCheckDelete(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "users.db")
    reg := bsureg.NewReg()
    err = reg.OpenDB(path)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    authModel := bsuser.NewAuth(reg)
    require.NoError(t, err)

    err = authModel.SeedUsers()
    require.NoError(t, err)

    contr := NewContr(authModel)

    auth := dsrpc.CreateAuth([]byte("admin"), []byte("admin"))

    addParams := bsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456"
    addResult := bsapi.NewAddUserResult()
    err = dsrpc.LocalExec(bsapi.AddUserMethod, addParams, addResult, auth, contr.AddUserHandler)
    require.NoError(t, err)

    addParams = bsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456xx"
    addResult = bsapi.NewAddUserResult()
    err = dsrpc.LocalExec(bsapi.AddUserMethod, addParams, addResult, auth, contr.AddUserHandler)
    require.Error(t, err)

    addParams = bsapi.NewAddUserParams()
    addParams.Login    = "йцукен"
    addParams.Pass     = "567890"
    addResult = bsapi.NewAddUserResult()
    err = dsrpc.LocalExec(bsapi.AddUserMethod, addParams, addResult, auth, contr.AddUserHandler)
    require.NoError(t, err)

    checkParams := bsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult := bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.NoError(t, err)
    require.Equal(t, true, checkResult.Match)

    checkParams = bsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456XXX"
    checkResult = bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.NoError(t, err)
    require.Equal(t, false, checkResult.Match)

    checkParams = bsapi.NewCheckUserParams()
    checkParams.Login    = "qwertyXXX"
    checkParams.Pass     = "123456"
    checkResult = bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.Error(t, err)
    require.Equal(t, false, checkResult.Match)

    deleteParams := bsapi.NewDeleteUserParams()
    deleteParams.Login    = "qwerty"
    deleteResult := bsapi.NewDeleteUserResult()
    err = dsrpc.LocalExec(bsapi.DeleteUserMethod, deleteParams, deleteResult, auth, contr.DeleteUserHandler)
    require.NoError(t, err)

    checkParams = bsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult = bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, auth, contr.CheckUserHandler)
    require.Error(t, err)
    require.Equal(t, false, checkResult.Match)

    err = reg.CloseDB()
    require.NoError(t, err)

}

func Test_User_Hello(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "blocks.db")
    reg := bsureg.NewReg()
    err = reg.OpenDB(path)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    store := bsuser.NewAuth(reg)
    require.NoError(t, err)

    err = store.SeedUsers()
    require.NoError(t, err)

    auth := dsrpc.CreateAuth([]byte("admin"), []byte("admin"))

    contr := NewContr(store)

    helloResp := GetHelloMsg
    params := bsapi.NewGetHelloParams()
    params.Message = GetHelloMsg
    result := bsapi.NewGetHelloResult()
    err = dsrpc.LocalExec(bsapi.GetHelloMethod, params, result, auth, contr.GetHelloHandler)

    require.NoError(t, err)
    require.Equal(t, helloResp, result.Message)

    err = reg.CloseDB()
    require.NoError(t, err)

}
