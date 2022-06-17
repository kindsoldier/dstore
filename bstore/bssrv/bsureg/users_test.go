package bsureg

import (
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
    "ndstore/bstore/bscom"
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

    login := "qwerty"

    descr1 := bscom.NewUserDescr()
    descr1.Login  = login
    descr1.Pass   = "12345"
    descr1.State  = "undef"
    descr1.Role   = "admin"

    err = reg.EraseUserDescr(login)
    require.NoError(t, err)

    err = reg.AddUserDescr(descr1)
    require.NoError(t, err)

    err = reg.AddUserDescr(descr1)
    require.Error(t, err)

    exists, err := reg.UserDescrExists(login)
    require.NoError(t, err)
    require.Equal(t, true, exists)

    descr2, err := reg.GetUserDescr(login)
    require.NoError(t, err)
    require.Equal(t, descr1.Login, descr2.Login)
    require.Equal(t, descr1.Pass, descr2.Pass)
    require.Equal(t, descr1.Role, descr2.Role)

    descr2.Pass = "56789"
    descr2.State = "disabled"
    descr2.Role = "somerole"

    err = reg.UpdateUserDescr(descr2)
    require.NoError(t, err)

    descr3, err := reg.GetUserDescr(login)
    require.NoError(t, err)
    require.NotNil(t, descr3)
    require.Equal(t, descr2.Login, descr3.Login)
    require.Equal(t, descr2.Pass, descr3.Pass)
    require.Equal(t, descr2.Role, descr3.Role)

    wrongLogin := "foobar"
    descr4, err := reg.GetUserDescr(wrongLogin)
    require.Error(t, err)
    require.NotNil(t, descr4)

    err = reg.EraseUserDescr(login)
    require.NoError(t, err)

    _, err = reg.GetUserDescr(login)
    require.Error(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}
