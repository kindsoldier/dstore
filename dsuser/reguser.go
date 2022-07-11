package dsuser

import (
    "strings"
    "dstore/dsdescr"
    "dstore/dsinter"
)

type Reg struct {
    db      dsinter.DB
    keyBase string
    sep     string
}

func NewReg(db dsinter.DB) (*Reg, error) {
    var err error
    var reg Reg
    reg.db      = db
    reg.keyBase = "user"
    reg.sep     = ":"
    return &reg, err
}

func (reg *Reg) Put(descr *dsdescr.User) error {
    var err error
    keyArr := []string{ reg.keyBase, descr.Login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) Has(login string) (bool, error) {
    var err error
    keyArr := []string{ reg.keyBase, login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) Get(login string) (*dsdescr.User, error) {
    var err error
    var descr *dsdescr.User
    keyArr := []string{ reg.keyBase, login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, err := reg.db.Get(keyBin)
    if err != nil {
        return descr, err
    }
    descr, err = dsdescr.UnpackUser(valBin)
    if err != nil {
        return descr, err
    }
    return descr, err
}

func (reg *Reg) Delete(login string) error {
    var err error
    keyArr := []string{ reg.keyBase, login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

//func(key []byte, val []byte) (bool, error)

func (reg *Reg) List() ([]*dsdescr.User, error) {
    var err error
    descrs := make([]*dsdescr.User, 0)
    cb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackUser(val)
        if err != nil {
            return interr, err
        }
        descrs = append(descrs, descr)
        return interr, err
    }
    err = reg.db.Iter(cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
