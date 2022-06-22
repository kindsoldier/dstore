/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "testing"
    "github.com/stretchr/testify/require"

    "ndstore/fstore/fssrv/fsreg"
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

    model := NewStore(rootDir, reg)
    require.NoError(t, err)

    err = model.SeedUsers()
    require.NoError(t, err)

    login   := "qwerty"
    pass    := "1234567"
    userName := "admin"

    err = model.DeleteUser(userName, login)
    require.NoError(t, err)

    err = model.AddUser(userName, login, pass)
    require.NoError(t, err)

    user, err := model.GetUser(login)
    require.NoError(t, err)
    require.Equal(t, login, user.Login)
    require.Equal(t, pass, user.Pass)

    ok, err := model.CheckUser(userName, login, pass)
    require.NoError(t, err)
    require.Equal(t, true, ok)

    err = reg.CloseDB()
    require.NoError(t, err)
}
