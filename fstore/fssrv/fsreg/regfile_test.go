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

func TestFile01(t *testing.T) {
    var err error
    var has bool

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)

    reg, err := NewReg(db)
    require.NoError(t, err)
    require.NotEqual(t, reg, nil)

    descr0 := dsdescr.NewFile()
    descr0.Login        = "admin"
    descr0.FilePath     = "/qwerty"
    descr0.FileId       = 2
    descr0.BatchCount   = 1
    descr0.DataSize     = 5
    descr0.CreatedAt    = 1657645101
    descr0.UpdatedAt    = 1657645102

    err = reg.PutFile(descr0)
    require.NoError(t, err)

    has, err = reg.HasFile(descr0.Login, descr0.FilePath)
    require.NoError(t, err)
    require.Equal(t, has, true)

    descr1, err := reg.GetFile(descr0.Login, descr0.FilePath)
    require.NoError(t, err)
    require.Equal(t, descr0, descr1)

    descrs, err := reg.ListFiles(descr0.Login)
    require.NoError(t, err)
    require.Equal(t, len(descrs), 1)

    err = reg.DeleteFile(descr0.Login, descr0.FilePath)
    require.NoError(t, err)

    has, err = reg.HasFile(descr0.Login, descr0.FilePath)
    require.NoError(t, err)
    require.Equal(t, has, false)
}
