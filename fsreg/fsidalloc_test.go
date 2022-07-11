/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsreg

import (
    "testing"
    "github.com/stretchr/testify/require"
    "dstore/dskeydb"
)

func TestIdAlloc(t *testing.T) {
    var err error

    dataDir := t.TempDir()
    db, err := dskeydb.OpenKV(dataDir, "tmp.leveldb")
    defer db.Close()
    require.NoError(t, err)

    key := []byte("blockids")
    alloc, err := NewAlloc(db, key)
    require.NoError(t, err)

    id1, err := alloc.NewId()
    require.NoError(t, err)

    err = alloc.FreeId(id1)
    require.NoError(t, err)

    id2, err := alloc.NewId()
    require.NoError(t, err)
    require.Equal(t, id1, id2)

    err = alloc.FreeId(id2)
    require.NoError(t, err)
}

func BenchmarkIdAlloc(b *testing.B) {
    var err error

    dataDir := b.TempDir()
    db := dskeydb.NewKV(dataDir, "tmp.leveldb")
    err = db.Open()
    defer db.Close()
    require.NoError(b, err)

    key := []byte("blockids")
    alloc, err := NewAlloc(db, key)
    require.NoError(b, err)

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            id, err := alloc.NewId()
            require.NoError(b, err)

            err = alloc.FreeId(id)
            require.NoError(b, err)
        }
    }
    b.SetParallelism(1000)
    b.RunParallel(pBench)
}
