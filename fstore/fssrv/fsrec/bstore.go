/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "ndstore/dscom"
)

func (store *Store) SeedBStores() error {
    var err error
    const address   = "127.0.0.1"
    var port      = "5100"
    const login     = "admin"
    const pass      = "admin"
    _, err = store.reg.AddBStoreDescr(address, port, login, pass, StateEnabled)
    port      = "5101"
    _, err = store.reg.AddBStoreDescr(address, port, login, pass, StateEnabled)
    port      = "5102"
    _, err = store.reg.AddBStoreDescr(address, port, login, pass, StateEnabled)
    return err
}

func (store *Store) AddBStore(address, port, login, pass string) error {
    var err error
    _, err = store.reg.AddBStoreDescr(address, port, login, pass, StateEnabled)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) GetBStore(address, port string) (*dscom.BStoreDescr, error) {
    var err error
    bStore, err := store.reg.GetBStoreDescr(address, port)
    return bStore, err
}

func (store *Store) UpdateBStore(address, port, login, pass string) error {
    var err error
    ok, err := checkPass(pass)
    if !ok {
        return err
    }
    err = store.reg.UpdateBStoreDescr(address, port, login, pass, StateEnabled)
    return err
}

func (store *Store) ListBStores() ([]*dscom.BStoreDescr, error) {
    var err error
    bStores, err := store.reg.ListBStoreDescrs()
    //for i := range BStores {
    //    BStores[i].Pass = "xxxxx"
    //}
    return bStores, err
}

func (store *Store) DeleteBStore(address, port string) error {
    var err error
    err = store.reg.DeleteBStoreDescr(address, port)
    return err
}
