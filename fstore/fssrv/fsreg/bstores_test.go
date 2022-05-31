package fsreg

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_BStoreDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dbPath := "postgres://pgsql@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    var address string  = "127.0.0.1"
    var port    string  = "5001"
    var login   string  = "qwerty"
    var pass    string  = "123456"
    var state   string  = "undef"

    err = reg.DeleteBStoreDescr(address, port)
    assert.NoError(t, err)

    id, err := reg.AddBStoreDescr(address, port, login, pass, state)
    assert.NoError(t, err)

    _, err = reg.AddBStoreDescr(address, port, login, pass, state)
    assert.Error(t, err)

    _, err = reg.AddBStoreDescr(address, port, login, pass, state)
    assert.Error(t, err)

    exists, err := reg.BStoreDescrExists(address, port)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    exists, err = reg.BStoreDescrExists(address + "xxx", port)
    assert.NoError(t, err)
    assert.Equal(t, false, exists)


    bstore, err := reg.GetBStoreDescr(address, port)
    assert.NoError(t, err)
    assert.Equal(t, id, bstore.BStoreId)
    assert.Equal(t, address, bstore.Address)
    assert.Equal(t, login, bstore.Login)
    assert.Equal(t, pass, bstore.Pass)

    login = "foobar"
    bstore.Login = login

    pass = "56789"
    bstore.Pass = pass

    state = "disabled"
    bstore.State = state

    err = reg.RenewBStoreDescr(bstore)
    assert.NoError(t, err)

    bstore, err = reg.GetBStoreDescr(address, port)
    assert.NoError(t, err)
    assert.Equal(t, id, bstore.BStoreId)
    assert.Equal(t, address, bstore.Address)
    assert.Equal(t, login, bstore.Login)
    assert.Equal(t, pass, bstore.Pass)

    err = reg.DeleteBStoreDescr(address, port)
    assert.NoError(t, err)

    _, err = reg.GetBStoreDescr(address, port)
    assert.Error(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
