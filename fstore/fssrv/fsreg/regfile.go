package fsreg

import (
    "strings"
    "dstore/dscomm/dsdescr"
)

func (reg *Reg) PutFile(descr *dsdescr.File) error {
    var err error
    keyArr := []string{ reg.fileBase, descr.Login, descr.FilePath }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasFile(login, filePath string) (bool, error) {
    var err error
    keyArr := []string{ reg.fileBase, login, filePath }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetFile(login, filePath string) (*dsdescr.File, error) {
    var err error
    var descr *dsdescr.File
    keyArr := []string{ reg.fileBase, login, filePath }
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

func (reg *Reg) DeleteFile(login, filePath string) error {
    var err error
    keyArr := []string{ reg.fileBase, login, filePath }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListFiles(login string) ([]*dsdescr.File, error) {
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


    keyArr := []string{ reg.fileBase, login }
    keyStr := strings.Join(keyArr, reg.sep)
    fileBaseBin := []byte(keyStr + reg.sep)
    err = reg.db.Iter(fileBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}

type FileFunc = func(fileDescr *dsdescr.File) (bool, error)

func (reg *Reg) ProcFiles(login string, fileCb FileFunc) ([]*dsdescr.File, error) {
    var err error
    descrs := make([]*dsdescr.File, 0)

    iterCb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackFile(val)
        if err != nil {
            return interr, err
        }
        interr, err = fileCb(descr)
        return interr, err
    }

    keyArr := []string{ reg.fileBase, login }
    keyStr := strings.Join(keyArr, reg.sep)
    fileBaseBin := []byte(keyStr + reg.sep)
    err = reg.db.Iter(fileBaseBin, iterCb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
