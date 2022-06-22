package fsfile

import (
    "bytes"
    "math/rand"
    "testing"
    "ndstore/fstore/fssrv/fsreg"
    "github.com/stretchr/testify/require"
)

func Test_File_WriteRead(t *testing.T) {
    var err error

    rootDir := t.TempDir()

    dbPath := "postgres://test@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    var batchSize   int64   = 5
    var blockSize   int64   = 1024

    fileId, file0, err := NewFile(reg, rootDir, batchSize, blockSize)
    require.NoError(t, err)
    require.NotEqual(t, file0, nil)

    var written int64
    var dataSize int64 = batchSize * blockSize * 10
    data := make([]byte, dataSize + 200)
    rand.Read(data)

    need := batchSize * blockSize * 2 + 10

    reader0 := bytes.NewReader(data)
    written, err = file0.Write(reader0, need)
    require.NoError(t, err)
    require.Equal(t, need, written)

    err = file0.Close()
    require.NoError(t, err)

    file3, err := OpenFile(reg, rootDir, fileId)
    require.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    read, err := file3.Read(writer)
    require.NoError(t, err)
    require.Equal(t, need, int64(len(writer.Bytes())))
    require.Equal(t, need, read)

    require.Equal(t, data[0:read], writer.Bytes())

    err = file3.Close()
    require.NoError(t, err)
}
