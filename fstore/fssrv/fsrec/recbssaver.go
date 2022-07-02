/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "fmt"
    "os"
    "ndstore/bstore/bsfunc"
    "ndstore/dscom"
    "ndstore/dsrpc"
    "ndstore/dserr"
    "ndstore/dslog"
)

type BSSaver struct {
    reg     dscom.IFSReg
    bstores []*dscom.BStoreDescr
    counter int
}

func NewBSSaver(reg dscom.IFSReg) *BSSaver {
    var saver BSSaver
    saver.reg    = reg
    saver.bstores = make([]*dscom.BStoreDescr, 0)
    return &saver
}

func (saver *BSSaver) LoadPool() error {
    var err error

    bstoresAll, err := saver.reg.ListBStoreDescrs()
    if err != nil {
        return dserr.Err(err)
    }

    bstores := make([]*dscom.BStoreDescr, 0)
    for _, bs := range bstoresAll {
        if bs.State == dscom.BStateNormal {
            bstores = append(bstores, bs)
        }
    }
    saver.bstores = bstores

    if len(saver.bstores) < 1 {
        err = errors.New("empty work store list")
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (saver *BSSaver) SaveBlock(descr *dscom.BlockDescr, blockReader *os.File) (int64, error) {
    var err error
    var bstoreId int64

    for {
        if len(saver.bstores) < 1 {
            err := errors.New("empty saver")
            return bstoreId, err
        }
        index := saver.counter % len(saver.bstores)

        bstore  := saver.bstores[index]
        uri     := fmt.Sprintf("%s:%s", bstore.Address, bstore.Port)
        login   := bstore.Login
        pass    := bstore.Pass
        auth    := dsrpc.CreateAuth([]byte(login), []byte(pass))
        bstoreId = bstore.BStoreId

        dslog.LogDebug("use bstoreId", bstoreId)

        blockReader.Seek(0, 0)
        binSize := descr.DataSize
        tryCount := 3
        for i := 0; i < tryCount; i++ {
            err = bsfunc.SaveBlock(uri, auth, descr, blockReader, binSize)
            if err == nil {
                break
            }
        }
        if err != nil {
            dslog.LogDebug("error save to bstoreId", bstoreId, uri)
            saver.deleteStore(index)
            continue
        }
        saver.counter++
        break
    }
    return bstoreId, dserr.Err(err)
}

func (saver *BSSaver) deleteStore(i int) {

    dslog.LogDebug("delete bstore index", i)

    stores := saver.bstores
    head := stores[0:i]
    tail := stores[(i + 1):len(stores)]
    stores = append(head, tail...)
    saver.bstores = stores
}
