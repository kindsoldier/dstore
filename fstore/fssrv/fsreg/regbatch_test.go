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

func Test_BatchDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    ver0 := time.Now().UnixNano()
    ver1 := time.Now().UnixNano()
    ver2 := time.Now().UnixNano()
    require.NotEqual(t, ver0, ver1)
    require.NotEqual(t, ver1, ver2)

    var fileId      int64   = 1
    var batchId     int64   = 2
    descr0 := dscom.NewBatchDescr()
    descr0.FileId   = fileId
    descr0.BatchId  = batchId
    descr0.BatchSize = 5
    descr0.BlockSize = 1024

    descr0.BatchVer = ver0
    descr0.UCounter = 1

    var exists bool

    _ = reg.EraseAllBatchDescrs()
    _ = reg.EraseSpecBatchDescr(descr0.FileId, descr0.BatchId, ver0)
    _ = reg.EraseSpecBatchDescr(descr0.FileId, descr0.BatchId, ver1)
    _ = reg.EraseSpecBatchDescr(descr0.FileId, descr0.BatchId, ver2)

    err = reg.AddNewBatchDescr(descr0)
    require.NoError(t, err)

    err = reg.AddNewBatchDescr(descr0)
    require.Error(t, err)

    descr0.BatchVer     = ver1
    err = reg.AddNewBatchDescr(descr0)
    require.NoError(t, err)

    descr0.BatchVer     = ver2
    err = reg.AddNewBatchDescr(descr0)
    require.NoError(t, err)

    descrs, err := reg.ListAllBatchDescrs()
    require.NoError(t, err)
    require.Equal(t, 3, len(descrs))

    exists, descr1, err := reg.GetNewestBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr0, descr1)


    err = reg.IncSpecBatchDescrUC(1, fileId, batchId, ver2)
    require.NoError(t, err)

    err = reg.DecSpecBatchDescrUC(1, fileId, batchId, ver2)
    require.NoError(t, err)

    err = reg.DecSpecBatchDescrUC(1, fileId, batchId, ver2)
    require.NoError(t, err)

    exists, descr2, err := reg.GetNewestBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr2.UCounter, int64(1))

    exists, descr3, err := reg.GetNewestBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr3.UCounter, int64(1))
    descr3.UCounter = descr0.UCounter
    descr3.BatchVer = descr0.BatchVer
    require.Equal(t, descr0, descr3)

    err = reg.EraseSpecBatchDescr(descr0.FileId, descr0.BatchId, ver0)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.EraseSpecBatchDescr(descr0.FileId, descr0.BatchId, ver1)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, false)

    err = reg.IncSpecBatchDescrUC(1, fileId, batchId, ver2)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, true)

    err = reg.EraseSpecBatchDescr(descr0.FileId, descr0.BatchId, ver2)
    require.NoError(t, err)

    exists, _, err = reg.GetNewestBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, false)
}
