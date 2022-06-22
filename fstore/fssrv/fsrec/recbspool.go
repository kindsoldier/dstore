/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "fmt"
    //"io"
    "os"
    //"math/rand"

    "ndstore/bstore/bsfunc"
    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dscom"
    "ndstore/dsrpc"
    "ndstore/dserr"
    "ndstore/dslog"
)


type BSPool struct {
    reg     *fsreg.Reg
    bstores []*dscom.BStoreDescr
    counter int
}

func NewBSPool(reg *fsreg.Reg) *BSPool {
    var pool BSPool
    pool.reg    = reg
    pool.bstores = make([]*dscom.BStoreDescr, 0)
    return &pool
}

func (pool *BSPool) LoadPool() error {
    var err error

    bstoresAll, err := pool.reg.ListBStoreDescrs()
    if err != nil {
        return dserr.Err(err)
    }

    bstores := make([]*dscom.BStoreDescr, 0)
    for _, bs := range bstoresAll {
        if bs.State == dscom.BStateNormal {
            bstores = append(bstores, bs)
        }
    }
    pool.bstores = bstores

    if len(pool.bstores) < 1 {
        err = errors.New("empty work store list")
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (pool *BSPool) deleteStore(i int) {
    dslog.LogDebug("delete bstore index", i)

    stores := pool.bstores
    head := stores[0:i]
    tail := stores[(i + 1):len(stores)]
    stores = append(head, tail...)
    pool.bstores = stores
}

func (pool *BSPool) SaveBlock(fileId, batchId, blockId, blockSize int64, blockReader *os.File, dataSize int64,
                                                blockType, hashAlg, hashInit, hashSum string) (int64, error) {
    var err error
    var bstoreId int64

    for {
        if len(pool.bstores) < 1 {
            err := errors.New("empty pool")
            return bstoreId, err
        }

        index := pool.counter % len(pool.bstores)

        bstore  := pool.bstores[index]

        uri     := fmt.Sprintf("%s:%s", bstore.Address, bstore.Port)
        login   := bstore.Login
        pass    := bstore.Pass
        auth    := dsrpc.CreateAuth([]byte(login), []byte(pass))
        bstoreId = bstore.BStoreId
        //dslog.LogDebug("use bstoreId", bstoreId)

        _ = bsfunc.DeleteBlock(uri, auth, fileId, batchId, blockId, blockType)
        //if err != nil {
        //    pool.deleteStore(index)
        //    continue
        //}
        blockReader.Seek(0, 0)
        err = bsfunc.SaveBlock(uri, auth, fileId, batchId, blockId, blockSize, blockReader,
                                                    dataSize, blockType, hashAlg, hashInit, hashSum)

        if err != nil {
            dslog.LogDebug("error save to bstoreId", bstoreId, bstore.Port)
            pool.deleteStore(index)
            continue
        }
        pool.counter++
        break
    }
    return bstoreId, dserr.Err(err)
}
