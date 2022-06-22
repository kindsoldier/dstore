/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "testing"
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

    var fileId      int64   = 1
    var batchId     int64   = 2
    descr0 := dscom.NewBatchDescr()
    descr0.FileId   = fileId
    descr0.BatchId  = batchId
    descr0.BatchSize = 5
    descr0.BlockSize = 1024

    var exists bool

    err = reg.AddBatchDescr(descr0)
    require.NoError(t, err)

    err = reg.AddBatchDescr(descr0)
    require.Error(t, err)

    descr7 := dscom.NewBatchDescr()
    *descr7 = *descr0
    descr7.FileId += 1
    require.NotEqual(t, descr7, descr0)
    err = reg.AddBatchDescr(descr7)
    require.NoError(t, err)

    exists, descr1, err := reg.GetBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr1, descr0)

    descrs, err := reg.ListBatchDescrs(fileId)
    require.NoError(t, err)
    require.Equal(t, len(descrs), 1)

    err = reg.EraseBatchDescr(fileId + 1, batchId)
    require.NoError(t, err)

    exists, descr2, err := reg.GetBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr2, descr0)

    err = reg.EraseBatchDescr(fileId, batchId)
    require.NoError(t, err)

    exists, _, err = reg.GetBatchDescr(fileId, batchId)
    require.NoError(t, err)
    require.Equal(t, exists, false)
}
