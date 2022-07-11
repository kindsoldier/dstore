/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsuser

import(
    "fmt"
    "testing"
    "github.com/stretchr/testify/require"

    "dstore/dskvdb"
    //"dstore/fsreg"
    //"dstore/dslog"
    "dstore/dsdescr"
)


func TestBlock01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.leveldb")
    defer db.Close()
    require.NoError(t, err)


    reg, err := NewReg(db)
    require.NoError(t, err)
    require.NotEqual(t, reg, nil)

    descr0 := dsdescr.NewUser()
    descr0.Login = "qwerty"
    descr0.Passw = "123456"
    descr0.CreatedAt = 8
    descr0.UpdatedAt = 9

    err = reg.Put(descr0)
    require.NoError(t, err)

    descr1, err := reg.Get(descr0.Login)
    require.NoError(t, err)
    require.Equal(t, descr0, descr1)

    descrs, err := reg.List()
    require.NoError(t, err)
    for _, descr := range descrs {
        descrBin, _ := descr.Pack()
        fmt.Println(string(descrBin))
    }
}
