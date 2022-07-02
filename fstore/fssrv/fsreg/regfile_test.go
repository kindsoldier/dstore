/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "testing"
    "time"
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
    ver0 := time.Now().UnixNano()
    ver1 := time.Now().UnixNano()
    ver2 := time.Now().UnixNano()
    require.NotEqual(t, ver0, ver1)
    require.NotEqual(t, ver1, ver2)

    descr0 := dscom.NewFileDescr()

    var fileId      int64   = 1
    descr0.FileId      = fileId
    descr0.BatchSize   = 5
    descr0.BlockSize   = 1024

    descr0.FileSize     = 5
    descr0.BlockSize    = 1024
    descr0.BatchCount   = 7

    descr0.FileVer  = ver0
    descr0.UCounter = 1

    var exists bool

    _ = reg.EraseAllFileDescrs()
    _ = reg.EraseSpecFileDescr(descr0.FileId, ver0)
    _ = reg.EraseSpecFileDescr(descr0.FileId, ver1)
    _ = reg.EraseSpecFileDescr(descr0.FileId, ver2)

    err = reg.AddNewFileDescr(descr0)
    require.NoError(t, err)

    err = reg.AddNewFileDescr(descr0)
    require.Error(t, err)

    descr0.FileVer     = ver1
    err = reg.AddNewFileDescr(descr0)
    require.NoError(t, err)

    descr0.FileVer     = ver2
    err = reg.AddNewFileDescr(descr0)
    require.NoError(t, err)

    descrs, err := reg.ListAllFileDescrs()
    require.NoError(t, err)
    require.Equal(t, 3, len(descrs))

    exists, descr1, err := reg.GetNewestFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr0, descr1)


    err = reg.IncSpecFileDescrUC(fileId, ver0)
    require.NoError(t, err)

    err = reg.DecSpecFileDescrUC(fileId, ver0)
    require.NoError(t, err)

    err = reg.DecSpecFileDescrUC(fileId, ver0)
    require.NoError(t, err)

    exists, descr2, err := reg.GetNewestFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr2.UCounter, int64(1))

    exists, descr3, err := reg.GetNewestFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr3.UCounter, int64(1))
    descr3.UCounter = descr0.UCounter
    descr3.FileVer = descr0.FileVer
    require.Equal(t, descr0, descr3)

    err = reg.EraseSpecFileDescr(descr0.FileId, ver0)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.EraseSpecFileDescr(descr0.FileId, ver1)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.IncSpecFileDescrUC(fileId, ver2)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.EraseSpecFileDescr(descr0.FileId, ver2)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestFileDescr(fileId)
    require.NoError(t, err)
    require.Equal(t, exists, false)

}
