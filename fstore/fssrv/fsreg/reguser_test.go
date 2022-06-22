package fsreg

import (
    "testing"
    "github.com/stretchr/testify/require"
    "ndstore/dscom"
)

func Test_UserDescr_InsertSelectDelete(t *testing.T) {
    var err error
    var exists bool

    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    var login   string  = "qwerty"

    descr0 := dscom.NewUserDescr()
    descr0.Login    = login
    descr0.Pass     = "123456"
    descr0.State    = "undef"
    descr0.Role     = "admin"

    err = reg.EraseUserDescr(login)
    require.NoError(t, err)

    _, err = reg.GetUserDescr(login)
    require.Error(t, err)

    id, err := reg.AddUserDescr(descr0)
    require.NoError(t, err)

    descr0.UserId = id

    id, err = reg.GetUserId(login)
    require.NoError(t, err)
    require.Equal(t, id, descr0.UserId)

    id, err = reg.GetUserId(login + "zzzz")
    require.Error(t, err)

    _, err = reg.AddUserDescr(descr0)
    require.Error(t, err)

    descr7 := dscom.NewUserDescr()
    *descr7 = *descr0
    descr7.Login += "xxxx"
    _, err = reg.AddUserDescr(descr7)
    require.NoError(t, err)
    require.NotEqual(t, descr7, descr0)

    _, err = reg.AddUserDescr(descr0)
    require.Error(t, err)

    exists, err = reg.UserDescrExists(login)
    require.NoError(t, err)
    require.Equal(t, true, exists)

    descr1, err := reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, descr0,     descr1)

    _, err = reg.GetUserDescr(login + "xxxx")
    require.NoError(t, err)

    descr1.Pass     = "56789"
    descr1.State    = "disabled"
    descr1.Role     = "somerole"
    err = reg.UpdateUserDescr(descr1)
    require.NoError(t, err)

    descr2, err := reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, descr2, descr1)

    descr2.Pass     = "567891XX"
    descr2.State    = "disabledXX"
    descr2.Role     = "someroleXX"
    err = reg.UpdateUserDescr(descr2)
    require.NoError(t, err)

    descr3, err := reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, descr2, descr3)

    err = reg.EraseUserDescr(login)
    require.NoError(t, err)

    _, err = reg.GetUserDescr(login)
    require.Error(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}
