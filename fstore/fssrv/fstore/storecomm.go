/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

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
}

func NewStore(dataDir string, reg dsinter.StoreReg) (*Store, error) {
    var err error
    var store Store
    store.dataDir   = dataDir
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
