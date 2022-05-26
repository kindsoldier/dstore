package fsreg

import (
    "math/rand"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestInsertSelectDeleteEntryDescr(t *testing.T) {
    var err error

    var dirPath     string = "x/y/z"
    var fileName    string = "qwerty.txt"
    var fileId      int64 = int64(rand.Intn(1024))

    dbPath := "postgres://pgsql@localhost/test"
    reg := NewReg()

    err = reg.OpenDB(dbPath)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    err = reg.DeleteEntryDescr(dirPath, fileName)
    assert.NoError(t, err)

    err = reg.AddEntryDescr(dirPath, fileName, fileId)
    assert.NoError(t, err)

    exists, err := reg.EntryDescrExists(dirPath, fileName)
    assert.NoError(t, err)
    assert.Equal(t, true, exists)

    entry, err := reg.GetEntryDescr(dirPath, fileName)
    assert.NoError(t, err)
    assert.Equal(t, fileId, entry.FileId)

    err = reg.DeleteEntryDescr(dirPath, fileName)
    assert.NoError(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
