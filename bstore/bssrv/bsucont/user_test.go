/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsucont

import (
    "testing"
    "path/filepath"

    "github.com/stretchr/testify/assert"

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
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    authModel := bsuser.NewAuth(reg)
    assert.NoError(t, err)

    contr := NewContr(authModel)

    addParams := bsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456"
    addResult := bsapi.NewAddUserResult()
    err = dsrpc.LocalExec(bsapi.AddUserMethod, addParams, addResult, nil, contr.AddUserHandler)
    assert.NoError(t, err)

    addParams = bsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456xx"
    addResult = bsapi.NewAddUserResult()
    err = dsrpc.LocalExec(bsapi.AddUserMethod, addParams, addResult, nil, contr.AddUserHandler)
    assert.Error(t, err)

    addParams = bsapi.NewAddUserParams()
    addParams.Login    = "йцукен"
    addParams.Pass     = "567890"
    addResult = bsapi.NewAddUserResult()
    err = dsrpc.LocalExec(bsapi.AddUserMethod, addParams, addResult, nil, contr.AddUserHandler)
    assert.NoError(t, err)

    checkParams := bsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult := bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.NoError(t, err)
    assert.Equal(t, true, checkResult.Match)

    checkParams = bsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456XXX"
    checkResult = bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.NoError(t, err)
    assert.Equal(t, false, checkResult.Match)

    checkParams = bsapi.NewCheckUserParams()
    checkParams.Login    = "qwertyXXX"
    checkParams.Pass     = "123456"
    checkResult = bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.Error(t, err)
    assert.Equal(t, false, checkResult.Match)

    deleteParams := bsapi.NewDeleteUserParams()
    deleteParams.Login    = "qwerty"
    deleteResult := bsapi.NewDeleteUserResult()
    err = dsrpc.LocalExec(bsapi.DeleteUserMethod, deleteParams, deleteResult, nil, contr.DeleteUserHandler)
    assert.NoError(t, err)

    checkParams = bsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult = bsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(bsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.Error(t, err)
    assert.Equal(t, false, checkResult.Match)

    err = reg.CloseDB()
    assert.NoError(t, err)

}

func Test_User_Hello(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "blocks.db")
    reg := bsureg.NewReg()
    err = reg.OpenDB(path)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    store := bsuser.NewAuth(reg)
    assert.NoError(t, err)

    contr := NewContr(store)

    helloResp := GetHelloMsg
    params := bsapi.NewGetHelloParams()
    params.Message = GetHelloMsg
    result := bsapi.NewGetHelloResult()
    err = dsrpc.LocalExec(bsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)

    err = reg.CloseDB()
    assert.NoError(t, err)

}
