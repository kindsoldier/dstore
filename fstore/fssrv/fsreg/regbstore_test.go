/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsreg

import(
    "testing"
    "github.com/stretchr/testify/require"

    "dstore/dscomm/dsdescr"
    "dstore/dscomm/dskvdb"
)

func TestBStore01(t *testing.T) {
    var err error
    var has bool

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)


    reg, err := NewReg(db)
    require.NoError(t, err)
    require.NotEqual(t, reg, nil)

    descr0 := dsdescr.NewBStore()
    descr0.Address   = "qwerty"
    descr0.Port      = "123456"
    descr0.CreatedAt = 1657645101
    descr0.UpdatedAt = 1657645102

    err = reg.PutBStore(descr0)
    require.NoError(t, err)

    has, err = reg.HasBStore(descr0.Address, descr0.Port)
    require.NoError(t, err)
    require.Equal(t, has, true)

    descr1, err := reg.GetBStore(descr0.Address, descr0.Port)
    require.NoError(t, err)
    require.Equal(t, descr0, descr1)

    descrs, err := reg.ListBStores()
    require.NoError(t, err)
    require.Equal(t, len(descrs), 1)

    err = reg.DeleteBStore(descr0.Address, descr0.Port)
    require.NoError(t, err)

    has, err = reg.HasBStore(descr0.Address, descr0.Port)
    require.NoError(t, err)
    require.Equal(t, has, false)
}
