package fsreg

import (
    "math/rand"
    "testing"

    "github.com/stretchr/testify/require"
)

func Test_EntryDescr_InsertSelectDelete(t *testing.T) {
    var err error

    var dirPath     string = "x/y/z"
    var fileName    string = "qwerty.txt"
    var fileId      int64 = int64(rand.Intn(1024))
    var userId      int64 = 8

    dbPath := "postgres://test@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    err = reg.EraseEntryDescr(userId, dirPath, fileName)
    require.NoError(t, err)

    err = reg.AddEntryDescr(userId, dirPath, fileName, fileId)
    require.NoError(t, err)

    exists, entry, err := reg.GetEntryDescr(userId, dirPath, fileName)
    require.NoError(t, err)
    require.Equal(t, fileId, entry.FileId)
    require.Equal(t, true, exists)

    err = reg.EraseEntryDescr(userId, dirPath, fileName)
    require.NoError(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}
