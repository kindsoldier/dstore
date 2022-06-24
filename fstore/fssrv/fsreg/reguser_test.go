package fsreg

import (
    "testing"
    "github.com/stretchr/testify/require"
    "ndstore/dscom"
)

func Test_UserDescr_InsertSelectDelete(t *testing.T) {
    var err error
    var exists bool

    // Open database
    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    // Create new user descriptor
    var login   string  = "qwerty"
    descr0 := dscom.NewUserDescr()
    descr0.Login    = login
    descr0.Pass     = "123456"
    descr0.State    = "undef"
    descr0.Role     = "admin"

    // Clean login id exist
    err = reg.EraseUserDescr(login)
    require.NoError(t, err)
    // Check user
    exists, _, err = reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, exists, false)

    // Add user
    id0, err := reg.AddUserDescr(descr0)
    require.NoError(t, err)
    descr0.UserId = id0
    // Again add the same user
    _, err = reg.AddUserDescr(descr0)
    require.Error(t, err)

    // Create new login
    newLogin := "xxxxxx"
    err = reg.EraseUserDescr(newLogin)
    require.NoError(t, err)
    // Add another user as copy prev user
    descr7 := dscom.NewUserDescr()
    *descr7 = *descr0
    descr7.Login = newLogin
    id7, err := reg.AddUserDescr(descr7)
    require.NoError(t, err)
    require.NotEqual(t, descr7, descr0)

    // Agan add user
    _, err = reg.AddUserDescr(descr7)
    require.Error(t, err)

    // Get user and compare data
    exists, descr8, err := reg.GetUserDescr(newLogin)
    require.NoError(t, err)
    descr7.UserId = id7
    require.Equal(t, descr7, descr8)
    require.Equal(t, exists, true)

    // Update user data
    descr0.Pass     = "56789"
    descr0.State    = "disabled"
    descr0.Role     = "somerole"
    err = reg.UpdateUserDescr(descr0)
    require.NoError(t, err)

    // Get user info and check
    exists, descr1, err := reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, exists, true)
    require.Equal(t, descr0, descr1)

    // Erase users and check
    err = reg.EraseUserDescr(login)
    require.NoError(t, err)

    exists, _, err = reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, exists, false)

    err = reg.EraseUserDescr(newLogin)
    require.NoError(t, err)

    exists, _, err = reg.GetUserDescr(newLogin)
    require.NoError(t, err)
    require.Equal(t, exists, false)

    // Close database
    err = reg.CloseDB()
    require.NoError(t, err)
}
