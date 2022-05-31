/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "testing"
    "github.com/stretchr/testify/assert"

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
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    model := fsrec.NewStore(rootDir, reg)
    assert.NoError(t, err)

    contr := NewContr(model)

    addParams := fsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456"
    addResult := fsapi.NewAddUserResult()
    err = dsrpc.LocalExec(fsapi.AddUserMethod, addParams, addResult, nil, contr.AddUserHandler)
    assert.NoError(t, err)

    checkParams := fsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult := fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.NoError(t, err)
    assert.Equal(t, true, checkResult.Match)


    addParams = fsapi.NewAddUserParams()
    addParams.Login    = "qwerty"
    addParams.Pass     = "123456xx"
    addResult = fsapi.NewAddUserResult()
    err = dsrpc.LocalExec(fsapi.AddUserMethod, addParams, addResult, nil, contr.AddUserHandler)
    assert.Error(t, err)

    addParams = fsapi.NewAddUserParams()
    addParams.Login    = "йцукен"
    addParams.Pass     = "567890"
    addResult = fsapi.NewAddUserResult()
    err = dsrpc.LocalExec(fsapi.AddUserMethod, addParams, addResult, nil, contr.AddUserHandler)
    assert.NoError(t, err)


    checkParams = fsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456XXX"
    checkResult = fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.NoError(t, err)
    assert.Equal(t, false, checkResult.Match)

    checkParams = fsapi.NewCheckUserParams()
    checkParams.Login    = "qwertyXXX"
    checkParams.Pass     = "123456"
    checkResult = fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.Error(t, err)
    assert.Equal(t, false, checkResult.Match)

    deleteParams := fsapi.NewDeleteUserParams()
    deleteParams.Login    = "qwerty"
    deleteResult := fsapi.NewDeleteUserResult()
    err = dsrpc.LocalExec(fsapi.DeleteUserMethod, deleteParams, deleteResult, nil, contr.DeleteUserHandler)
    assert.NoError(t, err)

    checkParams = fsapi.NewCheckUserParams()
    checkParams.Login    = "qwerty"
    checkParams.Pass     = "123456"
    checkResult = fsapi.NewCheckUserResult()
    err = dsrpc.LocalExec(fsapi.CheckUserMethod, checkParams, checkResult, nil, contr.CheckUserHandler)
    assert.Error(t, err)
    assert.Equal(t, false, checkResult.Match)

    err = reg.CloseDB()
    assert.NoError(t, err)

}
