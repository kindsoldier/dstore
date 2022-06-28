/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsbcont

import (
    "bytes"
    "math/rand"
    "testing"
    "path/filepath"

    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsblock"
    "ndstore/bstore/bssrv/bsbreg"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/require"
)




func Test_Block_SaveLoadDelete(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "tmp.blocks.db")
    reg := bsbreg.NewReg()
    err = reg.OpenDB(path)
    require.NoError(t, err)
    err = reg.MigrateDB()
    require.NoError(t, err)

    store := bsblock.NewStore(rootDir, reg)
    require.NoError(t, err)

    contr := NewContr(store)

    dataSize := int64(101)
    data := make([]byte, dataSize)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    params := bsapi.NewSaveBlockParams()
    params.FileId       = 2
    params.BatchId      = 3
    params.BlockId      = 4
    params.DataSize     = dataSize
    params.BlockSize    = 1024

    result := bsapi.NewSaveBlockResult()

    err = dsrpc.LocalPut(bsapi.SaveBlockMethod, reader, size, params, result, nil, contr.SaveBlockHandler)
    require.NoError(t, err)


    writer := bytes.NewBuffer(make([]byte, 0))

    err = dsrpc.LocalGet(bsapi.LoadBlockMethod, writer, params, result, nil, contr.LoadBlockHandler)
    require.NoError(t, err)
    require.Equal(t, len(data), len(writer.Bytes()))
    require.Equal(t, data, writer.Bytes())

    err = dsrpc.LocalExec(bsapi.DeleteBlockMethod, params, result, nil, contr.DeleteBlockHandler)
    require.NoError(t, err)

    err = reg.CloseDB()
    require.NoError(t, err)
}

func Test_Block_Hello(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "tmp.blocks.db")
    reg := bsbreg.NewReg()

    err = reg.OpenDB(path)
    require.NoError(t, err)

    err = reg.MigrateDB()
    require.NoError(t, err)

    store := bsblock.NewStore(rootDir, reg)
    require.NoError(t, err)

    contr := NewContr(store)

    helloResp := GetHelloMsg
    params := bsapi.NewGetHelloParams()
    params.Message = GetHelloMsg
    result := bsapi.NewGetHelloResult()
    err = dsrpc.LocalExec(bsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)

    require.NoError(t, err)
    require.Equal(t, helloResp, result.Message)

    err = reg.CloseDB()
    require.NoError(t, err)

}
