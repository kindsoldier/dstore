package fsreg

import (
    "strings"
    "strconv"
    "dstore/dsdescr"
)

func (reg *Reg) PutFile(descr *dsdescr.File) error {
    var err error
    idString := strconv.FormatInt(descr.FileId, 10)
    keyArr := []string{ reg.fileBase, idString }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasFile(fileId int64) (bool, error) {
    var err error
    idString := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.fileBase, idString }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetFile(fileId int64) (*dsdescr.File, error) {
    var err error
    var descr *dsdescr.File
    idString := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.fileBase, idString }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, err := reg.db.Get(keyBin)
    if err != nil {
        return descr, err
    }
    descr, err = dsdescr.UnpackFile(valBin)
    if err != nil {
        return descr, err
    }
    return descr, err
}

func (reg *Reg) DeleteFile(fileId int64) error {
    var err error
    idString := strconv.FormatInt(fileId, 10)
    keyArr := []string{ reg.fileBase, idString }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListFiles() ([]*dsdescr.File, error) {
    var err error
    descrs := make([]*dsdescr.File, 0)
    cb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackFile(val)
        if err != nil {
            return interr, err
        }
        descrs = append(descrs, descr)
        return interr, err
    }
    fileBaseBin := []byte(reg.fileBase + reg.sep)
    err = reg.db.Iter(fileBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
