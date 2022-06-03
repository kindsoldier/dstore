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
    var role    string  = "admin"

    err = reg.DeleteUserDescr(login)
    assert.NoError(t, err)

    err = reg.AddUserDescr(login, pass, state, role)
    assert.NoError(t, err)

    err = reg.AddUserDescr(login, pass, state, role)
    assert.Error(t, err)

    exists, err := reg.UserDescrExists(login)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    user, err := reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)
    assert.Equal(t, role, user.Role)

    pass = "56789"
    user.Pass = pass

    state = "disabled"
    user.State = state

    role = "somerole"
    user.Role = state

    err = reg.RenewUserDescr(user)
    assert.NoError(t, err)

    err = reg.UpdateUserDescr(login, pass, state, role)
    assert.NoError(t, err)

    user, err = reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)
    assert.Equal(t, role, user.Role)

    wrongLogin := "foobar"
    user, err = reg.GetUserDescr(wrongLogin)
    assert.Error(t, err)
    assert.NotNil(t, user)

    err = reg.DeleteUserDescr(login)
    assert.NoError(t, err)

    _, err = reg.GetUserDescr(login)
    assert.Error(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
