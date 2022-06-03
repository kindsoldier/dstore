/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "ndstore/fstore/fssrv/fsreg"
)

func Test_BStore_AddCheckDelete(t *testing.T) {
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

    err = model.SeedUsers()
    assert.NoError(t, err)

    address := "127.0.0.1"
    port    := "1234"
    login   := "qwerty"
    pass    := "123456"
    userName := "admin"
    err = model.AddBStore(userName, address, port, login, pass)
    assert.NoError(t, err)

    store, err := model.GetBStore(address, port)
    assert.NoError(t, err)
    assert.Equal(t, login, store.Login)
    assert.Equal(t, pass, store.Pass)

    err = model.DeleteBStore(userName, address, port)
    assert.NoError(t, err)

    store, err = model.GetBStore(address, port)
    assert.Error(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
