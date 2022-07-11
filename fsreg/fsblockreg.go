/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsreg

import (
    "fmt"
    "dstore/dsdescr"
    "dstore/dsinter"
)

type BlockReg struct {
    db      dsinter.KV
    keyBase []byte
}

func NewBlockReg(db dsinter.KV) (*BlockReg, error) {
    var err error
    var reg BlockReg
    reg.db      = db
    reg.keyBase = []byte("block")
    return &reg, err
}

func (reg *BlockReg) AddBlock(blockId int64, descr *dsdescr.Block) error {
    var err error

    key := fmt.Sprintf("%s.%20d", string(reg.keyBase), blockId)
    keyBin := []byte(key)

    has, err := reg.db.Has(keyBin)
    if err != nil {
        err = fmt.Errorf("add block %d error %s", blockId, err)
        return err
    }
    if has {
        err = fmt.Errorf("block %d yet exist", blockId)
        return err
    }
    descrBin, err := descr.Pack()
    err = reg.db.Put(keyBin, descrBin)
    if err != nil {
        err = fmt.Errorf("add block %d error %s", blockId, err)
        return err
    }
    return err
}

func (reg *BlockReg) GetBlock(blockId int64) (bool, *dsdescr.Block, error) {
    var err error
    var descr *dsdescr.Block

    key := fmt.Sprintf("%s.%20d", string(reg.keyBase), blockId)
    keyBin := []byte(key)

    has, err := reg.db.Has(keyBin)
    if err != nil {
        fmt.Errorf("get block %d error %s", blockId, err)
        return has, descr, err
    }
    if !has {
        return has, descr, err
    }

    descrBin, err := reg.db.Get(keyBin)
    if err != nil {
        fmt.Errorf("get block %d error %s", blockId, err)
        return has, descr, err
    }

    descr, err = dsdescr.UnpackBlock(descrBin)
    if err != nil {
        fmt.Errorf("unpack block %d error %s", blockId, err)
        return has, descr, err
    }
    return has, descr, err
}

func (reg *BlockReg) DeleteBlock(blockId int64) error {
    var err error

    key := fmt.Sprintf("%s.%20d", string(reg.keyBase), blockId)
    keyBin := []byte(key)

    has, err := reg.db.Has(keyBin)
    if err != nil {
        err = fmt.Errorf("del block %d error %s", blockId, err)
        return err
    }
    if !has {
        return err
    }
    err = reg.db.Delete(keyBin)
    if err != nil {
        err = fmt.Errorf("del block %d error %s", blockId, err)
        return err
    }
    return err
}
