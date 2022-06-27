package fsreg

import (
    "testing"
    "github.com/stretchr/testify/require"
    "ndstore/dscom"
)

func Test_BStoreDescr_InsertSelectErase(t *testing.T) {
    var err error

    // Create env
    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    var address string  = "127.0.0.1"
    var port    string  = "5001"
    // Create new descr
    descr0 := dscom.NewBStoreDescr()
    descr0.Address  = address
    descr0.Port     = port

    descr0.Login    = "qwerty"
    descr0.Pass     = "123456"
    descr0.State    = "undef"
    // Erase descr
    err = reg.EraseBStoreDescr(address, port)
    require.NoError(t, err)
    // Add descr
    id, err := reg.AddBStoreDescr(descr0)
    require.NoError(t, err)

    descr0.BStoreId = id
    // Add descr
    _, err = reg.AddBStoreDescr(descr0)
    require.Error(t, err)

    //descr7 := dscom.NewBStoreDescr()
    //*descr7 = *descr0
    //descr7.Port = "5007"
    //require.NotEqual(t, descr0, descr7)
    //_, err = reg.AddBStoreDescr(descr7)
    //require.NoError(t, err)

    // Get descr
    exists, descr1, err := reg.GetBStoreDescr(address, port)
    require.NoError(t, err)
    descr0.CreatedAt = 0
    descr0.UpdatedAt = 0
    descr1.CreatedAt = 0
    descr1.UpdatedAt = 0
    require.Equal(t, descr0, descr1)
    require.Equal(t, true, exists)
    // Update descr
    descr1.Login    = "foobar"
    descr1.Pass     = "56789"
    descr1.State    = "disabled"
    err = reg.UpdateBStoreDescr(descr1)
    require.NoError(t, err)

    exists, descr2, err := reg.GetBStoreDescr(address, port)
    require.NoError(t, err)
    descr2.CreatedAt = 0
    descr2.UpdatedAt = 0
    descr1.CreatedAt = 0
    descr1.UpdatedAt = 0
    require.Equal(t, descr1, descr2)
    require.Equal(t, true, exists)

    // Erase descr
    err = reg.EraseBStoreDescr(address, port)
    require.NoError(t, err)
    // Check existing
    exists, _, err = reg.GetBStoreDescr(address, port)
    require.NoError(t, err)
    require.Equal(t, false, exists)

    err = reg.CloseDB()
    require.NoError(t, err)
}
