/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsinter

import (
    "dstore/dsdescr"
)

type Tank interface {
     Write(data []byte) (int, error)
     Read(data []byte) (int, error)
     Clean() error
}

type IterFunc = func(key []byte, val []byte) (bool, error)

type DB interface {
    Put(key, val []byte) error
    Get(key []byte) ([]byte, error)
    Has(key []byte) (bool, error)
    Delete(key []byte) error
    Iter(cb IterFunc) error
}

type BlockReg interface {
    AddBlockDescr(descr *dsdescr.Block) error
    DelBlockDescr(blockId int64) error
}

type Alloc interface {
    NewId() (int64, error)
    FreeId(id int64) error
}
