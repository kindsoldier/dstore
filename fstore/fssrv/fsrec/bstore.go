/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "ndstore/dscom"
)

const BStateNormal      string = "normal"
const BStateDisabled    string = "disabled"
const BStateWrong       string = "wrong"

func (store *Store) SeedBStores() error {
    var err error

    bStores, err := store.reg.ListBStoreDescrs()
    if err != nil {
        return err
    }

    if len(bStores) > 0 {
        return err
    }

    const address   = "127.0.0.1"
    const login     = "admin"
    const pass      = "admin"
    ports := []string{ "5101", "5102", "5103" }
    for _, port := range ports {
        _, err = store.reg.AddBStoreDescr(address, port, login, pass, BStateNormal)
        if err != nil {
            return err
        }
    }
    return err
}

func (store *Store) AddBStore(userName, address, port, login, pass string) error {
    var err error

    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        return errors.New("insufficient rights for adding bStore")
    }

    _, err = store.reg.AddBStoreDescr(address, port, login, pass, BStateNormal)
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

func (store *Store) UpdateBStore(userName, address, port, login, pass string) error {
    var err error

    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        return errors.New("insufficient rights for updating bStore")
    }

    ok, err := validatePass(pass)
    if !ok {
        return err
    }
    err = store.reg.UpdateBStoreDescr(address, port, login, pass, BStateNormal)
    return err
}

func (store *Store) ListBStores(userName string) ([]*dscom.BStoreDescr, error) {
    var err error
    bStores := make([]*dscom.BStoreDescr, 0)

    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        return bStores, errors.New("insufficient rights for listing bStores")
    }

    bStores, err = store.reg.ListBStoreDescrs()
    //for i := range BStores {
    //    BStores[i].Pass = "xxxxx"
    //}
    if err != nil {
        return bStores, err
    }
    return bStores, err
}

func (store *Store) DeleteBStore(userName, address, port string) error {
    var err error
    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        return errors.New("insufficient rights for delete bStore")
    }
    err = store.reg.DeleteBStoreDescr(address, port)
    return err
}
