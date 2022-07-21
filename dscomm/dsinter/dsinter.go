/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsinter

import (
    "dstore/dscomm/dsdescr"
)

type IterFunc = func(key []byte, val []byte) (bool, error)
type DB interface {
    Put(key, val []byte) error
    Get(key []byte) ([]byte, error)
    Has(key []byte) (bool, error)
    Delete(key []byte) error
    Iter(prefix []byte, cb IterFunc) error
}

type Alloc interface {
    NewId() (int64, error)
    FreeId(id int64) error
    JSON() ([]byte, error)
    Syncer()
    Stop()
}

type Crate interface {
     Write(data []byte) (int, error)
     Read(data []byte) (int, error)
     Clean() error
}

type FStoreReg interface {
    PutUser(descr *dsdescr.User) error
    HasUser(login string) (bool, error)
    GetUser(login string) (*dsdescr.User, error)
    ListUsers() ([]*dsdescr.User, error)
    DeleteUser(login string) error

    DeleteFile(login, filePath string) error
    GetFile(login, filePath string) (*dsdescr.File, error)
    HasFile(login, filePath string) (bool, error)
    ListFiles(login string) ([]*dsdescr.File, error)
    PutFile(descr *dsdescr.File) error

    DeleteBatch(batchId, fileId int64) error
    GetBatch(batchId, fileId int64) (*dsdescr.Batch, error)
    HasBatch(batchId, fileId int64) (bool, error)
    ListBatchs() ([]*dsdescr.Batch, error)
    PutBatch(descr *dsdescr.Batch) error

    PutBlock(descr *dsdescr.Block) error
    GetBlock(fileId, batchId, blockType, blockId int64) (*dsdescr.Block, error)
    HasBlock(fileId, batchId, blockType, blockId int64) (bool, error)
    ListBlocks(fileId int64) ([]*dsdescr.Block, error)
    DeleteBlock(fileId, batchId, blockType, blockId int64) error


    DeleteBStore(address, port string) error
    GetBStore(address, port string) (*dsdescr.BStore, error)
    HasBStore(address, port string) (bool, error)
    ListBStores() ([]*dsdescr.BStore, error)
    PutBStore(descr *dsdescr.BStore) error
}

type BStoreReg interface {
    PutUser(descr *dsdescr.User) error
    HasUser(login string) (bool, error)
    GetUser(login string) (*dsdescr.User, error)
    ListUsers() ([]*dsdescr.User, error)
    DeleteUser(login string) error

    PutBlock(descr *dsdescr.Block) error
    GetBlock(fileId, batchId, blockType, blockId int64) (*dsdescr.Block, error)
    HasBlock(fileId, batchId, blockType, blockId int64) (bool, error)
    ListBlocks(fileId int64) ([]*dsdescr.Block, error)
    DeleteBlock(fileId, batchId, blockType, blockId int64) error
}
