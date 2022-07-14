package bsreg

import (
    "strings"
    "strconv"
    "dstore/dsdescr"
)

func (reg *Reg) PutBlock(descr *dsdescr.Block) error {
    var err error

    blockIdStr  := strconv.FormatInt(descr.BlockId, 10)
    batchIdStr  := strconv.FormatInt(descr.BatchId, 10)
    fileIdStr   := strconv.FormatInt(descr.FileId, 10)
    blockTypeStr := strconv.FormatInt(descr.BlockType, 10)

    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockTypeStr, blockIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasBlock(fileId, batchId, blockType, blockId int64) (bool, error) {
    var err error

    blockIdStr  := strconv.FormatInt(blockId, 10)
    batchIdStr  := strconv.FormatInt(batchId, 10)
    fileIdStr   := strconv.FormatInt(fileId, 10)
    blockTypeStr := strconv.FormatInt(blockType, 10)

    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockTypeStr, blockIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetBlock(fileId, batchId, blockType, blockId int64) (*dsdescr.Block, error) {
    var err error
    var descr *dsdescr.Block

    blockIdStr  := strconv.FormatInt(blockId, 10)
    batchIdStr  := strconv.FormatInt(batchId, 10)
    fileIdStr   := strconv.FormatInt(fileId, 10)
    blockTypeStr := strconv.FormatInt(blockType, 10)

    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockTypeStr, blockIdStr }
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

func (reg *Reg) DeleteBlock(fileId, batchId, blockType, blockId  int64) error {
    var err error

    blockIdStr  := strconv.FormatInt(blockId, 10)
    batchIdStr  := strconv.FormatInt(batchId, 10)
    fileIdStr   := strconv.FormatInt(fileId, 10)
    blockTypeStr := strconv.FormatInt(blockType, 10)

    keyArr := []string{ reg.blockBase, fileIdStr, batchIdStr, blockTypeStr, blockIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListBlocks(fileId int64) ([]*dsdescr.Block, error) {
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
    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.blockBase, fileIdStr }
    blockBaseStr := strings.Join(keyArr, reg.sep)
    blockBaseBin := []byte(blockBaseStr + reg.sep)
    err = reg.db.Iter(blockBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
