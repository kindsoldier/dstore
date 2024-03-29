/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsreg

import(
    "testing"
    "github.com/stretchr/testify/require"

    "dstore/dscomm/dsdescr"
    "dstore/dscomm/dskvdb"
)

func TestBlock01(t *testing.T) {
    var err error
    var has bool

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)


    reg, err := NewReg(db)
    require.NoError(t, err)
    require.NotEqual(t, reg, nil)

    descr0 := dsdescr.NewBlock()
    descr0.FileId     = 1
    descr0.BatchId    = 2
    descr0.BlockType  = 3
    descr0.BlockId    = 4
    descr0.BlockSize  = 1024
    descr0.DataSize   = 1001
    descr0.CreatedAt = 1657645101
    descr0.UpdatedAt = 1657645102

    err = reg.PutBlock(descr0)
    require.NoError(t, err)

    has, err = reg.HasBlock(descr0.FileId, descr0.BatchId, descr0.BlockType, descr0.BlockId )
    require.NoError(t, err)
    require.Equal(t, has, true)

    descr1, err := reg.GetBlock(descr0.FileId, descr0.BatchId, descr0.BlockType, descr0.BlockId )
    require.NoError(t, err)
    require.Equal(t, descr0, descr1)

    descrs, err := reg.ListBlocks(descr0.FileId)
    require.NoError(t, err)
    require.Equal(t, len(descrs), 1)

    err = reg.DeleteBlock(descr0.FileId, descr0.BatchId, descr0.BlockType, descr0.BlockId )
    require.NoError(t, err)

    has, err = reg.HasBlock(descr0.FileId, descr0.BatchId, descr0.BlockType, descr0.BlockId )
    require.NoError(t, err)
    require.Equal(t, has, false)
}
