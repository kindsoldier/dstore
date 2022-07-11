/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsreg

import (
    "sync"
    "dstore/dsdescr"
    "dstore/dsinter"
)

func NewBlockAlloc(db dsinter.KV) (*Alloc, error) {
    return NewAlloc(db, []byte("blockids"))
}

type Alloc struct {
    db      dsinter.KV
    topId   int64
    freeIds []int64
    key     []byte
    idMtx   sync.Mutex
}

func NewAlloc(db dsinter.KV, key []byte) (*Alloc, error) {
    var err error
    var alloc Alloc

    alloc.db        = db
    alloc.freeIds   = make([]int64, 0)
    alloc.key       = key
    alloc.topId     = 0

    has, err := alloc.db.Has(alloc.key)
    if err != nil {
        return &alloc, err
    }
    if !has {
        return &alloc, err
    }
    descrBin, err := alloc.db.Get(alloc.key)
    if err != nil {
        return &alloc, err
    }
    descr, err := dsdescr.UnpackAlloc(descrBin)
    if err != nil {
        return &alloc, err
    }
    alloc.freeIds = descr.FreeIds
    alloc.topId   = descr.TopId
    return &alloc, err
}

func (alloc *Alloc) NewId() (int64, error) {
    var err error
    var newId int64

    alloc.idMtx.Lock()
    defer alloc.idMtx.Unlock()

    freeIds := len(alloc.freeIds)
    if freeIds > 0 {
        newId = alloc.freeIds[freeIds - 1]
        alloc.freeIds = alloc.freeIds[0:freeIds - 1]
        return newId, err
    }
    alloc.topId++
    newId = alloc.topId

    err = alloc.storeState()
    if err != nil {
        alloc.freeIds = append(alloc.freeIds, newId)
        return newId, err
    }
    return newId, err
}

func (alloc *Alloc) FreeId(id int64) error {
    var err error

    alloc.idMtx.Lock()
    defer alloc.idMtx.Unlock()

    alloc.freeIds = append(alloc.freeIds, id)
    err = alloc.storeState()
    if err != nil {
        return err
    }
    return err
}

func (alloc *Alloc) toDescr() *dsdescr.Alloc {
    descr := dsdescr.NewAlloc()
    descr.TopId     = alloc.topId
    descr.FreeIds   = alloc.freeIds
    return descr
}

func (alloc *Alloc) storeState() error {
    var err error
    descr := alloc.toDescr()
    descrBin, err := descr.Pack()
    if err != nil {
        return err
    }
    err = alloc.db.Put(alloc.key, descrBin)
    if err != nil {
        return err
    }
    return err
}
