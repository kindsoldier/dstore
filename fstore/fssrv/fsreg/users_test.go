package fsreg

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func Test_UserDescr_InsertSelectDelete(t *testing.T) {
    var err error
    var exists bool

    dbPath := "postgres://pgsql@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    var login   string  = "qwerty"
    var pass    string  = "123456"
    var state   string  = "undef"
    var role    string  = "admin"

    err = reg.DeleteUserDescr(login)
    assert.NoError(t, err)

    _, err = reg.GetUserDescr(login)
    assert.Error(t, err)

    id, err := reg.AddUserDescr(login, pass, state, role)
    assert.NoError(t, err)

    userId, err := reg.GetUserId(login)
    assert.NoError(t, err)
    assert.Equal(t, id, userId)

    userId, err = reg.GetUserId(login + "zzzz")
    assert.Error(t, err)

    _, err = reg.AddUserDescr(login, pass, state, role)
    assert.Error(t, err)

    _, err = reg.AddUserDescr(login + "xxxx", pass, state, role)
    assert.NoError(t, err)

    _, err = reg.AddUserDescr(login, pass, state, role)
    assert.Error(t, err)

    exists, err = reg.UserDescrExists(login)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    user, err := reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.Equal(t, id, user.UserId)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)
    assert.Equal(t, state, user.State)
    assert.Equal(t, role, user.Role)

    _, err = reg.GetUserDescr(login + "xxxx")
    assert.NoError(t, err)

    pass = "56789"
    user.Pass = pass
    state = "disabled"
    user.State = state
    role = "somerole"
    user.Role = role

    err = reg.RenewUserDescr(user)
    assert.NoError(t, err)

    user, err = reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.Equal(t, id, user.UserId)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)
    assert.Equal(t, state, user.State)
    assert.Equal(t, role, user.Role)

    pass = "567891XX"
    user.Pass = pass
    state = "disabledXX"
    user.State = state
    role = "someroleXX"
    user.Role = role

    err = reg.UpdateUserDescr(login, pass, state, role)
    assert.NoError(t, err)

    user, err = reg.GetUserDescr(login)
    assert.NoError(t, err)
    assert.Equal(t, id, user.UserId)
    assert.Equal(t, login, user.Login)
    assert.Equal(t, pass, user.Pass)
    assert.Equal(t, state, user.State)
    assert.Equal(t, role, user.Role)

    err = reg.DeleteUserDescr(login)
    assert.NoError(t, err)

    _, err = reg.GetUserDescr(login)
    assert.Error(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
