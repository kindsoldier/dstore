/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "testing"
    "github.com/stretchr/testify/require"
    "ndstore/dscom"
)

func Test_FileDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)


    var exists bool
    descr0 := dscom.NewFileDescr()

    descr0.BatchSize   = 5
    descr0.BlockSize   = 1024
    descr0.UCounter    = 1
    descr0.BatchCount  = 7


    fileId, err := reg.AddFileDescr(descr0)
    require.NoError(t, err)

    descr0.FileId = fileId

    _, err = reg.AddFileDescr(descr0)
    require.NoError(t, err)

    descr7 := dscom.NewFileDescr()
    *descr7 = *descr0
    descr7.FileId += 1
    require.NotEqual(t, descr0, descr7)
    _, err = reg.AddFileDescr(descr7)
    require.NoError(t, err)

    descrs, err := reg.ListFileDescrs()
    require.NoError(t, err)
    require.GreaterOrEqual(t, len(descrs), 2)

    exists, descr1, err := reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr0, descr1)

    err = reg.EraseFileDescr(fileId + 1)
    require.NoError(t, err)

    exists, descr2, err := reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr0, descr2)

    descr2.BatchCount += 8
    err = reg.UpdateFileDescr(descr2)
    require.NoError(t, err)

    exists, descr3, err := reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr2, descr3)


    err = reg.EraseFileDescr(fileId)
    require.NoError(t, err)

    exists, _, err = reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, false)
}
