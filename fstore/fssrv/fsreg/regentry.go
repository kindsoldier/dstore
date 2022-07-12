package fsreg

import (
    "strings"
    "dstore/dsdescr"
)

func (reg *Reg) PutEntry(login string, descr *dsdescr.Entry) error {
    var err error
    keyArr := []string{ reg.entryBase, login, descr.FilePath }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasEntry(login, filePath string ) (bool, error) {
    var err error
    keyArr := []string{ reg.entryBase, login, filePath }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetEntry(login, filePath string) (*dsdescr.Entry, error) {
    var err error
    var descr *dsdescr.Entry
    keyArr := []string{ reg.entryBase, login, filePath }

    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, err := reg.db.Get(keyBin)
    if err != nil {
        return descr, err
    }
    descr, err = dsdescr.UnpackEntry(valBin)
    if err != nil {
        return descr, err
    }
    return descr, err
}

func (reg *Reg) DeleteEntry(login, filePath string) error {
    var err error
    keyArr := []string{ reg.entryBase, login, filePath }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListEntrys(login string) ([]*dsdescr.Entry, error) {
    var err error
    descrs := make([]*dsdescr.Entry, 0)

    cb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackEntry(val)
        if err != nil {
            return interr, err
        }
        descrs = append(descrs, descr)
        return interr, err
    }
    entryBaseArr := []string{ reg.entryBase, login }
    entryBaseBin := []byte(strings.Join(entryBaseArr, reg.sep) + reg.sep)
    err = reg.db.Iter(entryBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
