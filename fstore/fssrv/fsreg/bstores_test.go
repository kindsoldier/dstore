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

    var id      int64   = 1
    var address string  = "127.0.0.1:5001"
    var login   string  = "qwerty"
    var pass    string  = "123456"
    var state   string  = "undef"

    err = reg.DeleteBStoreDescr(id)
    assert.NoError(t, err)

    err = reg.AddBStoreDescr(id, address, login, pass, state)
    assert.NoError(t, err)

    err = reg.AddBStoreDescr(id, address, login, pass, state)
    assert.Error(t, err)

    exists, err := reg.BStoreDescrExists(id)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    bstore, _, err := reg.GetBStoreDescr(id)
    assert.NoError(t, err)
    assert.Equal(t, id, bstore.Id)
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

    bstore, _, err = reg.GetBStoreDescr(id)
    assert.NoError(t, err)
    assert.Equal(t, id, bstore.Id)
    assert.Equal(t, address, bstore.Address)
    assert.Equal(t, login, bstore.Login)
    assert.Equal(t, pass, bstore.Pass)

    err = reg.DeleteBStoreDescr(id)
    assert.NoError(t, err)

    _, exists, err = reg.GetBStoreDescr(id)
    assert.NoError(t, err)
    assert.Equal(t, false, exists)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
