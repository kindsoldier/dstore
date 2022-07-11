/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsreg

import (
    "testing"
    "math/rand"
    "time"

    "github.com/stretchr/testify/require"

    "dstore/dsdescr"
    "dstore/dskeydb"
)

func TestBlockReg(t *testing.T) {
    var err error
    var has bool

    dataDir := t.TempDir()
    db, err := dskeydb.OpenKV(dataDir, "tmp.leveldb")
    defer db.Close()
    require.NoError(t, err)

    reg, err := NewBlockReg(db)
    require.NoError(t, err)

    blockId := int64(1)

    descr0 := dsdescr.NewBlock()
    descr0.BlockSize = 1024
    descr0.DataSize  = 1000
    require.NoError(t, err)

    err = reg.AddBlock(blockId, descr0)
    require.NoError(t, err)

    err = reg.AddBlock(blockId, descr0)
    require.Error(t, err)

    has, descr1, err := reg.GetBlock(1)
    require.NoError(t, err)
    require.Equal(t, has, true)
    require.Equal(t, descr0, descr1)

    err = reg.DeleteBlock(blockId)
    require.NoError(t, err)

    err = reg.DeleteBlock(blockId)
    require.NoError(t, err)

    has, _, err = reg.GetBlock(blockId)
    require.NoError(t, err)
    require.Equal(t, has, false)
}


func BenchmarkBlockReg(b *testing.B) {
    var err error

    dataDir := b.TempDir()
    db, err := dskeydb.OpenKV(dataDir, "tmp.leveldb")
    defer db.Close()
    require.NoError(b, err)

    reg, err := NewBlockReg(db)
    require.NoError(b, err)

    alloc, err := NewBlockAlloc(db)
    require.NoError(b, err)

    rand.Seed(time.Now().UnixNano())

    pBench := func(pb *testing.PB) {
        for pb.Next() {

            blockId, err := alloc.NewId()
            require.NoError(b, err)

            descr0 := dsdescr.NewBlock()
            descr0.BlockSize = 1024
            descr0.DataSize  = 1000

            err = reg.AddBlock(blockId, descr0)
            require.NoError(b, err)

            has, descr1, err := reg.GetBlock(blockId)
            require.NoError(b, err)
            require.Equal(b, has, true)
            require.Equal(b, descr0, descr1)

            err = reg.DeleteBlock(blockId)
            require.NoError(b, err)

            alloc.FreeId(blockId)
        }
    }
    b.SetParallelism(1000)
    b.RunParallel(pBench)
}
