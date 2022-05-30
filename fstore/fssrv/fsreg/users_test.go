package fsreg

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func Test_UserDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dbPath := "postgres://pgsql@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)


    var id      int64   = 1
    var login   string  = "qwerty"
    var pass    string  = "123456"
    var state   string  = "undef"

    err = reg.DeleteUserDescr(id)
    assert.NoError(t, err)

    err = reg.AddUserDescr(id, login, pass, state)
    assert.NoError(t, err)

    err = reg.AddUserDescr(id, login, pass, state)
    assert.Error(t, err)

    exists, err := reg.UserDescrExists(id)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    user, _, err := reg.GetUserDescr(id)
    assert.NoError(t, err)
    assert.Equal(t, id, user.Id)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)

    login = "foobar"
    user.Login = login

    pass = "56789"
    user.Pass = pass

    state = "disabled"
    user.State = state

    err = reg.RenewUserDescr(user)
    assert.NoError(t, err)

    user, _, err = reg.GetUserDescr(id)
    assert.NoError(t, err)
    assert.Equal(t, id, user.Id)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)

    err = reg.DeleteUserDescr(id)
    assert.NoError(t, err)

    _, exists, err = reg.GetUserDescr(id)
    assert.NoError(t, err)
    assert.Equal(t, false, exists)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
