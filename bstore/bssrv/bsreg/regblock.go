package bsreg

import (
    "strings"
    "strconv"
    "dstore/dsdescr"
)

func (reg *Reg) PutBlock(descr *dsdescr.Block) error {
    var err error
    blockIdStr := strconv.FormatInt(descr.BlockId, 10)
    batchIdStr := strconv.FormatInt(descr.BatchId, 10)
    fileIdStr := strconv.FormatInt(descr.FileId, 10)
    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasBlock(blockId, batchId, fileId int64) (bool, error) {
    var err error
    blockIdStr := strconv.FormatInt(blockId, 10)
    batchIdStr := strconv.FormatInt(batchId, 10)
    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetBlock(blockId, batchId, fileId int64) (*dsdescr.Block, error) {
    var err error
    var descr *dsdescr.Block
    blockIdStr := strconv.FormatInt(blockId, 10)
    batchIdStr := strconv.FormatInt(batchId, 10)
    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, err := reg.db.Get(keyBin)
    if err != nil {
        return descr, err
    }
    descr, err = dsdescr.UnpackBlock(valBin)
    if err != nil {
        return descr, err
    }
    return descr, err
}

func (reg *Reg) DeleteBlock(blockId, batchId, fileId int64) error {
    var err error
    blockIdStr := strconv.FormatInt(blockId, 10)
    batchIdStr := strconv.FormatInt(batchId, 10)
    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListBlocks() ([]*dsdescr.Block, error) {
    var err error
    descrs := make([]*dsdescr.Block, 0)
    cb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackBlock(val)
        if err != nil {
            return interr, err
        }
        descrs = append(descrs, descr)
        return interr, err
    }
    blockBaseBin := []byte(reg.blockBase + reg.sep)
    err = reg.db.Iter(blockBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
