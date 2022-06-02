package fsreg

import (
    "math/rand"
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_EntryDescr_InsertSelectDelete(t *testing.T) {
    var err error

    var dirPath     string = "x/y/z"
    var fileName    string = "qwerty.txt"
    var fileId      int64 = int64(rand.Intn(1024))
    var userId      int64 = 8

    dbPath := "postgres://pgsql@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    err = reg.DeleteEntryDescr(userId, dirPath, fileName)
    assert.NoError(t, err)

    err = reg.AddEntryDescr(userId, dirPath, fileName, fileId)
    assert.NoError(t, err)

    exists, err := reg.EntryDescrExists(userId, dirPath, fileName)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    entry, err := reg.GetEntryDescr(userId, dirPath, fileName)
    assert.NoError(t, err)
    assert.Equal(t, fileId, entry.FileId)

    err = reg.DeleteEntryDescr(userId, dirPath, fileName)
    assert.NoError(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
