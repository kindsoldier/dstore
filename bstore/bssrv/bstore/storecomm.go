/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bstore

import (
    "io/fs"
    "time"
    "dstore/dsinter"
)

type Store struct {
    dataDir     string
    reg         dsinter.StoreReg
    dirPerm     fs.FileMode
    filePerm    fs.FileMode
    startTime   int64

    fileAlloc   dsinter.Alloc
}

func NewStore(dataDir string, reg dsinter.StoreReg, fileAlloc dsinter.Alloc) (*Store, error) {
    var err error
    var store Store
    store.dataDir   = dataDir
    store.fileAlloc = fileAlloc
    store.reg       = reg
    store.dirPerm   = 0755
    store.filePerm  = 0644
    store.startTime = time.Now().Unix()
    return &store, err
}

func (store *Store) SetDirPerm(dirPerm fs.FileMode) {
    store.dirPerm = dirPerm
}

func (store *Store) SetFilePerm(filePerm fs.FileMode) {
    store.filePerm = filePerm
}

func (store *Store) GetUptime() (int64, error) {
    var err error
    uptime := time.Now().Unix() - store.startTime
    return uptime, err
}
