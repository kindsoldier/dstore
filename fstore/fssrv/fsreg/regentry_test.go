/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import(
    "testing"
    "github.com/stretchr/testify/require"

    "dstore/dsdescr"
    "dstore/dskvdb"
)

func TestEntry01(t *testing.T) {
    var err error
    var has bool

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)


    reg, err := NewReg(db)
    require.NoError(t, err)
    require.NotEqual(t, reg, nil)

    descr0 := dsdescr.NewEntry()
    descr0.FilePath = "qwerty"
    descr0.FileId   = 1
    descr0.CreatedAt = 1657645101
    descr0.UpdatedAt = 1657645102

    login := "admin"

    err = reg.PutEntry(login, descr0)
    require.NoError(t, err)

    has, err = reg.HasEntry(login, descr0.FilePath)
    require.NoError(t, err)
    require.Equal(t, has, true)

    descr1, err := reg.GetEntry(login, descr0.FilePath)
    require.NoError(t, err)
    require.Equal(t, descr0, descr1)

    descrs, err := reg.ListEntrys(login)
    require.NoError(t, err)
    require.Equal(t, len(descrs), 1)

    err = reg.DeleteEntry(login, descr0.FilePath)
    require.NoError(t, err)

    has, err = reg.HasEntry(login, descr0.FilePath)
    require.NoError(t, err)
    require.Equal(t, has, false)
}
