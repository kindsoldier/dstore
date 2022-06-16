package bsureg

import (
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
)

func Test_UserDescr_InsertSelectDelete(t *testing.T) {
    var err error

    dataRoot := t.TempDir()
    path := filepath.Join(dataRoot, "users.db")
    reg := NewReg()
    err = reg.OpenDB(path)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    var login   string  = "qwerty"
    var pass    string  = "12345"
    var state   string  = "undef"
    var role    string  = "admin"

    err = reg.DeleteUserDescr(login)
    require.NoError(t, err)

    err = reg.AddUserDescr(login, pass, state, role)
    require.NoError(t, err)

    err = reg.AddUserDescr(login, pass, state, role)
    require.Error(t, err)

    exists, err := reg.UserDescrExists(login)
    require.NoError(t, err)
    require.Equal(t, true, exists)

    user, err := reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, login, user.Login)
    require.Equal(t, pass, user.Pass)
    require.Equal(t, role, user.Role)

    pass = "56789"
    user.Pass = pass

    state = "disabled"
    user.State = state

    role = "somerole"
    user.Role = state

    err = reg.RenewUserDescr(user)
    require.NoError(t, err)

    err = reg.UpdateUserDescr(login, pass, state, role)
    require.NoError(t, err)

    user, err = reg.GetUserDescr(login)
    require.NoError(t, err)
    require.NotNil(t, user)
    require.Equal(t, login, user.Login)
    require.Equal(t, pass, user.Pass)
    require.Equal(t, role, user.Role)

    wrongLogin := "foobar"
    user, err = reg.GetUserDescr(wrongLogin)
    require.Error(t, err)
    require.NotNil(t, user)

    err = reg.DeleteUserDescr(login)
    require.NoError(t, err)

    _, err = reg.GetUserDescr(login)
    require.Error(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}
