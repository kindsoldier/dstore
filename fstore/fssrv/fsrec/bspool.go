/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "fmt"
    "io"
    "ndstore/bstore/bsfunc"
    "ndstore/fstore/fssrv/fsreg"
    "ndstore/dscom"
    "ndstore/dsrpc"
    "ndstore/dserr"
)


type BSPool struct {
    reg     *fsreg.Reg
    bstores []*dscom.BStoreDescr
    counter int64
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


func (pool *BSPool) SaveBlock(fileId, batchId, blockId, blockSize int64, blockReader io.Reader, dataSize int64,
                                                blockType, hashAlg, hashInit, hashSum string) (int64, error) {
    var err error
    var bstoreId int64

    //index := pool.counter % int64(len(pool.bstores))

    bstore := pool.bstores[0]

    uri     := fmt.Sprintf("%s:%s", bstore.Address, bstore.Port)
    login   := bstore.Login
    pass    := bstore.Pass
    auth    := dsrpc.CreateAuth([]byte(login), []byte(pass))

    err = bsfunc.SaveBlock(uri, auth, fileId, batchId, blockId, blockSize, blockReader,
                                                dataSize, blockType, hashAlg, hashInit, hashSum)
    if err != nil {
        return bstoreId, dserr.Err(err)
    }

    pool.counter++

    return bstoreId, dserr.Err(err)
}
