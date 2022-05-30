/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "io/fs"
    "ndstore/fstore/fssrv/fsreg"
)

const blockFileExt string = ".blk"
const storeDBName  string = "file.db"


type Store struct {
    dataRoot string
    dirPerm   fs.FileMode
    filePerm  fs.FileMode
    reg    *fsreg.Reg
}

func NewStore(dataRoot string, reg *fsreg.Reg) *Store {
    var store Store
    store.dataRoot  = dataRoot
    store.dirPerm   = 0755
    store.filePerm  = 0644
    store.reg       = reg
    return &store
}

func (store *Store) SetPerm(dirPerm, filePerm fs.FileMode) {
    store.dirPerm = dirPerm
    store.filePerm = filePerm
}
