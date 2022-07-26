package fsreg

import (
    "strings"
    "strconv"
    "dstore/dscomm/dsdescr"
)

func (reg *Reg) PutBatch(descr *dsdescr.Batch) error {
    var err error
    batchIdStr := strconv.FormatInt(descr.BatchId, 10)
    fileIdStr := strconv.FormatInt(descr.FileId, 10)
    keyArr := []string{ reg.batchBase, fileIdStr, batchIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasBatch(fileId, batchId int64) (bool, error) {
    var err error
    batchIdStr := strconv.FormatInt(batchId, 10)
    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.batchBase, fileIdStr, batchIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetBatch(fileId, batchId int64) (*dsdescr.Batch, error) {
    var err error
    var descr *dsdescr.Batch
    batchIdStr := strconv.FormatInt(batchId, 10)
    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.batchBase, fileIdStr, batchIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, err := reg.db.Get(keyBin)
    if err != nil {
        return descr, err
    }
    descr, err = dsdescr.UnpackBatch(valBin)
    if err != nil {
        return descr, err
    }
    return descr, err
}

func (reg *Reg) DeleteBatch(fileId, batchId int64) error {
    var err error
    batchIdStr := strconv.FormatInt(batchId, 10)
    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.batchBase, fileIdStr, batchIdStr }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListBatchs(fileId int64) ([]*dsdescr.Batch, error) {
    var err error
    descrs := make([]*dsdescr.Batch, 0)
    cb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackBatch(val)
        if err != nil {
            return interr, err
        }
        descrs = append(descrs, descr)
        return interr, err
    }

    fileIdStr := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.batchBase, fileIdStr }
    batchBaseStr := strings.Join(keyArr, reg.sep)
    batchBaseBin := []byte(batchBaseStr + reg.sep)
    err = reg.db.Iter(batchBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
