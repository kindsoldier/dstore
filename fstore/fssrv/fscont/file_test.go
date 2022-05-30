package fdcont

import (
    "bytes"
    "math/rand"
    "testing"

    "github.com/stretchr/testify/assert"

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
    assert.NoError(t, err)
    err = reg.MigrateDB()
    assert.NoError(t, err)

    store := fsrec.NewStore(t.TempDir(), reg)
    contr := NewContr(store)


    fileName := "../aaaa/qwert.txt"
    saveParams := fsapi.NewSaveFileParams()
    saveParams.FilePath = fileName
    saveResult := fsapi.NewSaveFileResult()

    data := make([]byte, 8)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    err = dsrpc.LocalPut(fsapi.SaveFileMethod, reader, size, saveParams, saveResult, nil, contr.SaveFileHandler)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    loadParams := fsapi.NewLoadFileParams()
    loadParams.FilePath = fileName
    loadResult := fsapi.NewLoadFileResult()

    err = dsrpc.LocalGet(fsapi.LoadFileMethod, writer, loadParams, loadResult, nil, contr.LoadFileHandler)
    assert.NoError(t, err)
    assert.Equal(t, len(data), len(writer.Bytes()))
    assert.Equal(t, data, writer.Bytes())

    listParams := fsapi.NewListFilesParams()
    listParams.DirPath = "../../aaaa/"
    listResult := fsapi.NewListFilesResult()

    err = dsrpc.LocalExec(fsapi.ListFilesMethod, listParams, listResult, nil, contr.ListFilesHandler)
    assert.NoError(t, err)

    for _, entry := range listResult.Entries {
        fmt.Println(entry.DirPath, entry.FileName)
    }
    //deleteParams := fsapi.NewDeleteFileParams()
    //deleteParams.FilePath = "qwert.txt"
    //deleteResult := fsapi.NewDeleteFileResult()
    //err = dsrpc.LocalExec(fsapi.DeleteFileMethod, deleteParams, deleteResult, nil, contr.DeleteFileHandler)
    //assert.NoError(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}
