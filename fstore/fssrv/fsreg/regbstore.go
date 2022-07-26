package fsreg

import (
    "strings"
    "dstore/dscomm/dsdescr"
)

func (reg *Reg) PutBStore(descr *dsdescr.BStore) error {
    var err error
    keyArr := []string{ reg.bstoreBase, descr.Address, descr.Port }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasBStore(address, port string) (bool, error) {
    var err error
    keyArr := []string{ reg.bstoreBase, address, port }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetBStore(address, port string) (*dsdescr.BStore, error) {
    var err error
    var descr *dsdescr.BStore
    keyArr := []string{ reg.bstoreBase, address, port }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, err := reg.db.Get(keyBin)
    if err != nil {
        return descr, err
    }
    descr, err = dsdescr.UnpackBStore(valBin)
    if err != nil {
        return descr, err
    }
    return descr, err
}

func (reg *Reg) DeleteBStore(address, port string) error {
    var err error
    keyArr := []string{ reg.bstoreBase, address, port }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListBStores() ([]*dsdescr.BStore, error) {
    var err error
    descrs := make([]*dsdescr.BStore, 0)
    cb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackBStore(val)
        if err != nil {
            return interr, err
        }
        descrs = append(descrs, descr)
        return interr, err
    }
    bstoreKeyBaseBin := []byte(reg.bstoreBase + reg.sep)
    err = reg.db.Iter(bstoreKeyBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
