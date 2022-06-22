package fdcont

import (
    "bytes"
    "math/rand"
    "testing"

    "github.com/stretchr/testify/require"

    "ndstore/fstore/fsapi"
    "ndstore/fstore/fssrv/fsrec"
    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dsrpc"
    "fmt"

)


func Test_File_SaveLoadDelete(t *testing.T) {
    var err error

    dbPath := "postgres://pgsql@localhost/test"
    reg := fsreg.NewReg()

    err = reg.OpenDB(dbPath)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    store := fsrec.NewStore(t.TempDir(), reg)
    contr := NewContr(store)

    err = store.SeedUsers()
    require.NoError(t, err)

    auth := dsrpc.CreateAuth([]byte("admin"), []byte("admin"))

    fileName := "../aaaa/qwert.txt"
    saveParams := fsapi.NewSaveFileParams()
    saveParams.FilePath = fileName
    saveResult := fsapi.NewSaveFileResult()

    data := make([]byte, 8)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dsrpc.LocalPut(fsapi.SaveFileMethod, reader, size, saveParams, saveResult, auth, contr.SaveFileHandler)
    require.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    loadParams := fsapi.NewLoadFileParams()
    loadParams.FilePath = fileName
    loadResult := fsapi.NewLoadFileResult()

    err = dsrpc.LocalGet(fsapi.LoadFileMethod, writer, loadParams, loadResult, auth, contr.LoadFileHandler)
    require.NoError(t, err)
    require.Equal(t, len(data), len(writer.Bytes()))
    require.Equal(t, data, writer.Bytes())

    listParams := fsapi.NewListFilesParams()
    listParams.DirPath = "../../aaaa/"
    listResult := fsapi.NewListFilesResult()

    err = dsrpc.LocalExec(fsapi.ListFilesMethod, listParams, listResult, auth, contr.ListFilesHandler)
    require.NoError(t, err)

    for _, entry := range listResult.Entries {
        fmt.Println(entry.DirPath, entry.FileName)
    }
    //deleteParams := fsapi.NewDeleteFileParams()
    //deleteParams.FilePath = "qwert.txt"
    //deleteResult := fsapi.NewDeleteFileResult()
    //err = dsrpc.LocalExec(fsapi.DeleteFileMethod, deleteParams, deleteResult, auth, contr.DeleteFileHandler)
    //require.NoError(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}
