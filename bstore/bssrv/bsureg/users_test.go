package bsureg

import (
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
)

func Test_UserDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dataRoot := t.TempDir()
    path := filepath.Join(dataRoot, "users.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    var login   string  = "qwerty"
    var pass    string  = "12345"
    var state   string  = "undef"

    err = reg.DeleteUserDescr(login)
    assert.NoError(t, err)

    err = reg.AddUserDescr(login, pass, state)
    assert.NoError(t, err)

    err = reg.AddUserDescr(login, pass, state)
    assert.Error(t, err)

    exists, err := reg.UserDescrExists(login)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    user, _, err := reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)

    pass = "56789"
    user.Pass = pass

    state = "disabled"
    user.State = state

    err = reg.RenewUserDescr(user)
    assert.NoError(t, err)

    err = reg.UpdateUserDescr(login, pass, state)
    assert.NoError(t, err)

    user, _, err = reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)

    wrongLogin := "foobar"
    user, _, err = reg.GetUserDescr(wrongLogin)
    assert.NoError(t, err)
    assert.Nil(t, user)

    err = reg.DeleteUserDescr(login)
    assert.NoError(t, err)

    _, exists, err = reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.Equal(t, false, exists)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
