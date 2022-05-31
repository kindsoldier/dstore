/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "testing"
    "github.com/stretchr/testify/assert"

    "ndstore/fstore/fssrv/fsreg"
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

    model := NewStore(rootDir, reg)
    assert.NoError(t, err)


    login := "qwerty"
    pass := "1234567"
    err = model.AddUser(login, pass)
    assert.NoError(t, err)

    user, err := model.GetUser(login)
    assert.NoError(t, err)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)

    ok, err := model.CheckUser(login, pass)
    assert.NoError(t, err)
    assert.Equal(t, true, ok)


    err = reg.CloseDB()
    assert.NoError(t, err)
}
