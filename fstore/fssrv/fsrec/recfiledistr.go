/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "fmt"
    "bytes"
    "math/rand"
    "ndstore/bstore/bsfunc"
    "ndstore/dscom"
    "ndstore/dsrpc"
    "ndstore/dserr"
    "ndstore/dslog"
    "ndstore/fstore/fssrv/fsfile"
)

type FileDistr struct {
    reg     dscom.IFSReg
    dataDir string
    bstores []*dscom.BStoreDescr
    counter int
}

func NewFileDistr(dataDir string, reg dscom.IFSReg) *FileDistr {
    var distr FileDistr
    distr.dataDir = dataDir
    distr.reg = reg
    distr.bstores = make([]*dscom.BStoreDescr, 0)
    return &distr
}

func (distr *FileDistr) LoadPool() error {
    var err error

    bstoresAll, err := distr.reg.ListBStoreDescrs()
    if err != nil {
        return dserr.Err(err)
    }

    bstores := make([]*dscom.BStoreDescr, 0)
    for _, bs := range bstoresAll {
        if bs.State == dscom.BStateNormal {
            bstores = append(bstores, bs)
        }
    }
    distr.bstores = bstores

    distr.counter = rand.Intn(len(distr.bstores))

    if len(distr.bstores) < 1 {
        err = errors.New("empty work store list")
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (distr *FileDistr) SaveBlock(descr *dscom.BlockDescr) (bool, int64, error) {
    var err error
    var bstoreId int64
    ok := false

    for {
        if len(distr.bstores) < 1 {
            err := errors.New("empty distr")
            return ok, bstoreId, err
        }
        index := distr.counter % len(distr.bstores)

        bstore  := distr.bstores[index]
        uri     := fmt.Sprintf("%s:%s", bstore.Address, bstore.Port)
        login   := bstore.Login
        pass    := bstore.Pass
        auth    := dsrpc.CreateAuth([]byte(login), []byte(pass))
        bstoreId = bstore.BStoreId

        block, err := fsfile.OpenBlock(distr.reg, distr.dataDir,  descr.FileId, descr.BatchId, descr.BlockId, descr.BlockType)
        defer block.Close()

        writer := bytes.NewBuffer(make([]byte, 0))
        _, err = block.Read(writer)
        if err != nil {
            dslog.LogDebug("error save to bstoreId", bstoreId, uri)
            return ok, bstoreId, err
        }

        binSize := descr.DataSize

        tryCount := 3
        for i := 0; i < tryCount; i++ {
            reader := bytes.NewReader(writer.Bytes())
            err = bsfunc.SaveBlock(uri, auth, descr, reader, binSize)
            if err == nil {
                break
            }
        }
        if err != nil {
            dslog.LogDebug("error save to bstoreId", bstoreId, uri)
            distr.deleteStore(index)
            continue
        }
        distr.counter++
        ok = true
        break
    }
    return ok, bstoreId, dserr.Err(err)
}

func (distr *FileDistr) deleteStore(i int) {

    dslog.LogDebug("delete bstore index", i)

    stores := distr.bstores
    head := stores[0:i]
    tail := stores[(i + 1):len(stores)]
    stores = append(head, tail...)
    distr.bstores = stores
}
