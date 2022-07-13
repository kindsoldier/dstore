/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsfile

import(
    "bytes"
    "math/rand"
    "testing"
    "github.com/stretchr/testify/require"

    "dstore/dskvdb"
    "dstore/fstore/fssrv/fsreg"
)

func TestFile01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.leveldb")
    defer db.Close()
    require.NoError(t, err)

    reg, err := fsreg.NewReg(db)
    require.NoError(t, err)

    filePath    := "/qwerty.txt"
    login       := "admin"
    var fileId      int64 = 3
    var batchSize   int64 = 5
    var blockSize   int64 = 1000 * 1000
    var batchCount  int64 = 10

    file, err := NewFile(reg, dataDir, login, filePath,  fileId, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, file, nil)

    dataSize := batchCount * batchSize * blockSize
    origin := make([]byte, dataSize)
    rand.Read(origin)
    reader := bytes.NewReader(origin)

    needSize := int64(dataSize)
    wrSize, err := file.Write(reader, needSize)
    require.NoError(t, err)
    require.Equal(t, needSize, wrSize)

    _, err = reg.GetFile(login, filePath)
    require.NoError(t, err)

    file, err = OpenFile(reg, dataDir, login, filePath)
    require.NoError(t, err)
    require.NotEqual(t, file, nil)

    writer := bytes.NewBuffer(nil)

    readSize, err := file.Read(writer)
    require.NoError(t, err)
    require.Equal(t, wrSize, readSize)
    //require.Equal(t, origin[0:wrSize], writer.Bytes())

    err = file.Clean()
    require.NoError(t, err)
}
