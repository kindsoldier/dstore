/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsalloc

import (
    "sync"
    "context"
    "time"

    encoder "github.com/vmihailenco/msgpack/v5"

    "dstore/dscomm/dsinter"
    "dstore/dscomm/dslog"
)

type Alloc struct {
    db          dsinter.DB
    topId       int64
    freeIds     []int64
    key         []byte
    clean       bool
    giantMtx    sync.Mutex

    ctx     context.Context
    cancel  context.CancelFunc
    wg      sync.WaitGroup
}

func OpenAlloc(db dsinter.DB, key []byte) (*Alloc, error) {
    var err error
    var alloc Alloc

    alloc.db        = db
    alloc.freeIds   = make([]int64, 0)
    alloc.key       = key
    alloc.topId     = 0

    alloc.ctx, alloc.cancel = context.WithCancel(context.Background())

    has, err := alloc.db.Has(alloc.key)
    if err != nil {
        return &alloc, err
    }
    if has {
        descrBin, err := alloc.db.Get(alloc.key)
        if err != nil {
            return &alloc, err
        }
        descr, err := UnpackAllocDescr(descrBin)
        if err != nil {
            return &alloc, err
        }
        if !descr.Clean {
            // todo: ?
        }
        alloc.freeIds  = descr.FreeIds
        alloc.topId    = descr.TopId
        alloc.clean    = false
    }
    return &alloc, err
}

func (alloc *Alloc) NewId() (int64, error) {
    var err error
    var newId int64

    alloc.giantMtx.Lock()
    defer alloc.giantMtx.Unlock()

    freeIds := len(alloc.freeIds)
    if freeIds > 0 {
        newId = alloc.freeIds[freeIds - 1]
        alloc.freeIds = alloc.freeIds[0:freeIds - 1]
        return newId, err
    }

    newId = alloc.topId + 1
    alloc.topId = newId
    return newId, err
}

func (alloc *Alloc) FreeId(id int64) error {
    var err error

    alloc.giantMtx.Lock()
    defer alloc.giantMtx.Unlock()

    switch {
        case id == alloc.topId:
            alloc.topId--
        case id > alloc.topId:  // todo: ???
            return err
        default:
            alloc.freeIds = append(alloc.freeIds, id)
    }
    return err
}

func (alloc *Alloc) JSON() ([]byte, error) {
    var err error
    descr := alloc.toDescr()
    descrBin, err := descr.Pack()
    if err != nil {
        return descrBin, err
    }
    return descrBin, err
}


func (alloc *Alloc) toDescr() *AllocDescr {
    descr := NewAllocDescr()
    alloc.giantMtx.Lock()
    descr.TopId     = alloc.topId
    descr.FreeIds   = alloc.freeIds
    descr.Clean     = alloc.clean
    alloc.giantMtx.Unlock()
    return descr
}


func (alloc *Alloc) Stop() {
    alloc.cancel()
    alloc.wg.Wait()
}

func (alloc *Alloc) Syncer()  {
    alloc.wg.Add(1)
    lastDescr := alloc.toDescr()
    for {
        time.Sleep(1000 * time.Millisecond)
        select {
            case <- alloc.ctx.Done():
                alloc.saveState()
                dslog.LogInfo("alloc loop canceled")
                alloc.wg.Done()
                return
            default:
        }
        descr := alloc.toDescr()
        if descr.TopId != lastDescr.TopId || len(descr.FreeIds) != len(lastDescr.FreeIds) {
            begin := time.Now()
            descrBin, err := descr.Pack()
            if err != nil {
                dslog.LogErrorf("alloc loop pack error: %v", err)
                continue
            }
            err = alloc.db.Put(alloc.key, descrBin)
            if err != nil {
                dslog.LogErrorf("alloc loop put error: %v", err)
                continue
            }
            used := time.Since(begin)
            dslog.LogDebugf("alloc saving time: %v", used)
        }
        lastDescr = descr
    }
}

func (alloc *Alloc) saveState() error  {
    var err error
    descr := alloc.toDescr()
    alloc.giantMtx.Lock()           // todo: its rapid
    descr.Clean = true
    descrBin, err := descr.Pack()
    if err != nil {
        dslog.LogErrorf("alloc loop pack error: %v", err)
        return err
    }
    err = alloc.db.Put(alloc.key, descrBin)
    if err != nil {
        dslog.LogErrorf("alloc loop put error: %v", err)
        return err
    }
    return err
}

type AllocDescr struct {
    TopId   int64           `json:"topId"	msgpack:"topId"`
    FreeIds []int64         `json:"freeIds"	msgpack:"freeIds"`
    Clean   bool            `json:"clean"	msgpack:"clean"`
}

func NewAllocDescr() *AllocDescr {
    var descr AllocDescr
    return &descr
}

func UnpackAllocDescr(descrBin []byte) (*AllocDescr, error) {
    var err error
    var descr AllocDescr
    err = encoder.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *AllocDescr) Pack() ([]byte, error) {
    var err error
    descrBin, err := encoder.Marshal(descr)
    return descrBin, err
}
