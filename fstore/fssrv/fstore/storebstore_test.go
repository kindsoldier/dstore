/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "testing"
    "github.com/stretchr/testify/require"

    "dstore/dscomm/dskvdb"
    "dstore/dscomm/dsdescr"
    "dstore/fstore/fssrv/fsreg"
)


func TestBStore01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "storedb")
    defer db.Close()
    require.NoError(t, err)

    reg, err := fsreg.NewReg(db)
    require.NoError(t, err)

    store, err := NewStore(dataDir, reg, nil)
    require.NoError(t, err)

    err = store.SeedUsers()
    require.NoError(t, err)

    //err = store.SeedBStores()
    //require.NoError(t, err)

    descr0 := dsdescr.NewBStore()
    descr0.Address    = "localhost"
    descr0.Port       = "1234"
    descr0.Login      = "admin"
    descr0.Pass       = "admin"

    adminLogin   := "admin"
    wrongLogin   := "wrong"

    err = store.AddBStore(wrongLogin, descr0)
    require.Error(t, err)

    err = store.AddBStore(adminLogin, descr0)
    require.NoError(t, err)

    descrs, err := store.ListBStores(adminLogin, "")
    require.NoError(t, err)
    require.Equal(t, len(descrs), 1)

    has, descr1, err := store.GetBStore(descr0.Address, descr0.Port)
    require.NoError(t, err)
    require.Equal(t, has, true)
    require.Equal(t, descr0, descr1)

    var ok bool
    ok, err = store.CheckBStore(adminLogin, descr0.Address, descr0.Port, descr0.Login, descr0.Pass)
    require.NoError(t, err)
    require.Equal(t, true, ok)

    ok, err = store.CheckBStore(wrongLogin, descr0.Address, descr0.Port, descr0.Login, descr0.Pass)
    require.Error(t, err)
    require.Equal(t, false, ok)

    err = store.DeleteBStore(wrongLogin, descr0.Address, descr0.Port)
    require.Error(t, err)

    err = store.DeleteBStore(adminLogin, descr0.Address, descr0.Port)
    require.NoError(t, err)

    err = store.DeleteBStore(descr0.Login, descr0.Address, descr0.Port)
    require.NoError(t, err)

    _, err = store.ListBStores(wrongLogin, "")
    require.Error(t, err)

    descrs, err = store.ListBStores(adminLogin, "")
    require.NoError(t, err)
    require.Equal(t, len(descrs), 0)
}
