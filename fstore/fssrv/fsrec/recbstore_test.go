/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "testing"
    "github.com/stretchr/testify/require"
    "ndstore/fstore/fssrv/fsreg"
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

    model := NewStore(rootDir, reg)
    require.NoError(t, err)

    err = model.SeedUsers()
    require.NoError(t, err)

    address := "127.0.0.1"
    port    := "1234"
    login   := "qwerty"
    pass    := "123456"
    userName := "admin"
    err = model.AddBStore(userName, address, port, login, pass)
    require.NoError(t, err)

    store, err := model.GetBStore(address, port)
    require.NoError(t, err)
    require.Equal(t, login, store.Login)
    require.Equal(t, pass, store.Pass)

    err = model.DeleteBStore(userName, address, port)
    require.NoError(t, err)

    store, err = model.GetBStore(address, port)
    require.Error(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}
