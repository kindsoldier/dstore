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

    // Create env
    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)
    err = reg.MigrateDB()
    require.NoError(t, err)

    // Create origins
    var exists bool
    descr0 := dscom.NewFileDescr()

    descr0.BatchSize   = 5
    descr0.BlockSize   = 1024
    descr0.UCounter    = 1
    descr0.BatchCount  = 7

    // Add file descr
    fileId, err := reg.AddFileDescr(descr0)
    require.NoError(t, err)

    descr0.FileId = fileId
    // Again add the same file descr
    _, err = reg.AddFileDescr(descr0)
    require.NoError(t, err)

    // Add secondary file descr
    descr7 := dscom.NewFileDescr()
    *descr7 = *descr0
    descr7.FileId += 1
    require.NotEqual(t, descr0, descr7)
    _, err = reg.AddFileDescr(descr7)
    require.NoError(t, err)

    // Lists file descrs
    descrs, err := reg.ListFileDescrs()
    require.NoError(t, err)
    require.GreaterOrEqual(t, len(descrs), 2)

    // Get file descrs
    exists, descr1, err := reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    // Compare dates
    require.NotEqual(t, descr0.CreatedAt, descr1.CreatedAt)
    require.NotEqual(t, descr0.UpdatedAt, descr1.UpdatedAt)
    // Zeroind dates for compare descrs
    descr0.CreatedAt = 0
    descr1.CreatedAt = 0
    descr0.UpdatedAt = 0
    descr1.UpdatedAt = 0
    require.Equal(t, descr0, descr1)

    // Erase notexists desc
    err = reg.EraseFileDescr(fileId + 1)
    require.NoError(t, err)

    // Get file descr
    exists, descr2, err := reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.NotEqual(t, descr0.CreatedAt, descr2.CreatedAt)
    require.NotEqual(t, descr0.UpdatedAt, descr2.UpdatedAt)
    descr0.CreatedAt = 0
    descr2.CreatedAt = 0
    descr0.UpdatedAt = 0
    descr2.UpdatedAt = 0
    require.Equal(t, descr0, descr2)

    // Update file descrs
    descr2.BatchCount += 8
    err = reg.UpdateFileDescr(descr2)
    require.NoError(t, err)

    // Again get file descr
    exists, descr3, err := reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.NotEqual(t, descr2.CreatedAt, descr3.CreatedAt)
    require.NotEqual(t, descr2.UpdatedAt, descr3.UpdatedAt)
    descr3.CreatedAt = 0
    descr2.CreatedAt = 0
    descr3.UpdatedAt = 0
    descr2.UpdatedAt = 0

    require.Equal(t, descr2, descr3)

    // Erase file descr
    err = reg.EraseFileDescr(fileId)
    require.NoError(t, err)

    // Get erased file descrs
    exists, _, err = reg.GetFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, false)
}
