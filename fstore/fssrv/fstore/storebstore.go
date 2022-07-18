/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fstore

import (
    "errors"
    "fmt"
    "regexp"
    "time"
    "dstore/dscomm/dsdescr"
    "dstore/dscomm/dserr"
)

func (store *Store) SeedBStores() error {
    var err error

    bStores, err := store.reg.ListBStores()
    if err != nil {
        return dserr.Err(err)
    }
    if len(bStores) > 0 {
        return dserr.Err(err)
    }
    ports := []string{ "5101", "5102", "5103", "5104", "5105", "5106", "5107" }
    descr := dsdescr.NewBStore()
    descr.Address  = "127.0.0.1"
    descr.Login    = "admin"
    descr.Pass     = "admin"
    descr.State    = dsdescr.BSStateEnabled
    descr.CreatedAt = time.Now().Unix()
    descr.UpdatedAt = descr.CreatedAt

    for _, port := range ports {
        descr.Port = port
        err = store.reg.PutBStore(descr)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (store *Store) AddBStore(login string, bstore *dsdescr.BStore) error {
    var err error
    var ok bool

    role, err := store.getUserRole(login)
    if role != dsdescr.URoleAdmin {
        err = fmt.Errorf("insufficient rights for %s", login)
        return dserr.Err(err)
    }


    ok, err = validateBSAddess(bstore.Address)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateBSPort(bstore.Port)
    if !ok {
        return dserr.Err(err)
    }

    has, err := store.reg.HasBStore(bstore.Address, bstore.Port)
    if err != nil {
        return dserr.Err(err)
    }
    if has {
        err = fmt.Errorf("address:port %s:%s exist", bstore.Address, bstore.Port)
        return dserr.Err(err)
    }
    bstore.State  = dsdescr.BSStateEnabled
    bstore.CreatedAt = time.Now().Unix()
    bstore.UpdatedAt = bstore.CreatedAt

    err = store.reg.PutBStore(bstore)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) GetBStore(login, port string) (bool, *dsdescr.BStore, error) {
    var err error
    var bstore *dsdescr.BStore
    has, err := store.reg.HasBStore(login, port)
    if err != nil {
        return has, bstore, dserr.Err(err)
    }
    if !has {
        return has, bstore, dserr.Err(err)
    }
    bstore, err = store.reg.GetBStore(login, port)
    if err != nil {
        return has, bstore, dserr.Err(err)
    }
    return has,bstore, dserr.Err(err)
}

func (store *Store) CheckBStore(authLogin, address, port, login, pass string) (bool, error) {
    var err error
    var ok, ok1, ok2 bool

    userRole, err := store.getUserRole(authLogin)
    if userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", authLogin)
        return ok, dserr.Err(err)
    }

    has, err := store.reg.HasBStore(address, port)
    if err != nil {
        return ok, dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("bstore %s not exist", login)
    }
    descr, err := store.reg.GetBStore(address, port)
    if err != nil {
        return ok, dserr.Err(err)
    }
    if login == descr.Login {
        ok1 = true
    }

    if pass == descr.Pass {
        ok2 = true
    }
    ok = ok1 && ok2
    return ok, dserr.Err(err)
}

func (store *Store) UpdateBStore(login string, bstore *dsdescr.BStore) error {
    var err error
    // Get current role
    userRole, err := store.getUserRole(login)
    if err != nil {
        return dserr.Err(err)
    }
    // Set defaults
    // Rigth control
    if userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", login)
        return dserr.Err(err)
    }

    // Get old profile and copy to new
    oldBStore, err := store.reg.GetBStore(bstore.Address, bstore.Port)
    if err != nil {
        return dserr.Err(err)
    }
    newBStore := dsdescr.NewBStore()
    newBStore.Address     = oldBStore.Address
    newBStore.Port        = oldBStore.Port
    newBStore.Login       = oldBStore.Login
    newBStore.Pass        = oldBStore.Pass
    newBStore.State       = oldBStore.State
    newBStore.CreatedAt   = oldBStore.CreatedAt
    newBStore.UpdatedAt   = time.Now().Unix()

    // Update property if exists
    if len(bstore.Address) > 0 {
        newBStore.Address = bstore.Address
    }
    if len(bstore.Port) > 0 {
        newBStore.Port = bstore.Port
    }
    if len(bstore.Login) > 0 {
        newBStore.Login = bstore.Login
    }
    if len(bstore.Pass) > 0 {
        newBStore.Pass = bstore.Pass
    }
    if len(bstore.State) > 0 {
        newBStore.State = bstore.State
    }

    // Validation new property
    var ok bool
    ok, err = validateBSAddess(newBStore.Address)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateBSPort(newBStore.Port)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateBSState(newBStore.State)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateBSLogin(newBStore.Login)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateBSPass(newBStore.Pass)
    if !ok {
        return dserr.Err(err)
    }
    // Delete old bstore descr
    err = store.reg.DeleteBStore(bstore.Address, bstore.Port)
    if err != nil {
        return dserr.Err(err)
    }
    // Put new bstore descr
    err = store.reg.PutBStore(newBStore)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) ListBStores(authLogin, regular string) ([]*dsdescr.BStore, error) {
    var err error

    resDescrs := make([]*dsdescr.BStore, 0)
    userRole, err := store.getUserRole(authLogin)
    if userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", authLogin)
        return resDescrs, dserr.Err(err)
    }

    descrs, err := store.reg.ListBStores()
    if err != nil {
        return resDescrs, dserr.Err(err)
    }
    if len(regular) == 0 {
        resDescrs = descrs
        return resDescrs, dserr.Err(err)
    }
    re, err := regexp.CompilePOSIX(regular)
    if err != nil {
        return resDescrs, dserr.Err(err)
    }
    for _, descr := range descrs {
        ok := re.Match([]byte(descr.Address))
        if ok {
            resDescrs = append(resDescrs, descr)
        }
    }
    return resDescrs, dserr.Err(err)
}

func (store *Store) DeleteBStore(login string, address, port string) error {
    var err error

    userRole, err := store.getUserRole(login)
    if userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", login)
        return dserr.Err(err)
    }

    err = store.reg.DeleteBStore(address, port)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func validateBSAddess(address string) (bool, error) {
    var err error
    var ok bool = true
    if len(address) == 0 {
        ok = false
        err = errors.New("zero len address")
    }
    return ok, dserr.Err(err)
}

func validateBSPort(port string) (bool, error) {
    var err error
    var ok bool = true
    if len(port) == 0 {
        ok = false
        err = errors.New("zero len address")
    }
    return ok, dserr.Err(err)
}

func validateBSState(state string) (bool, error) {
    var err error
    var ok bool = true
    if state == dsdescr.BSStateDisabled  {
        return ok, dserr.Err(err)
    }
    if state == dsdescr.BSStateEnabled  {
        return ok, dserr.Err(err)
    }
    err = errors.New("irrelevant state name")
    ok = false
    return ok, dserr.Err(err)
}

func validateBSLogin(login string) (bool, error) {
    var err error
    var ok bool = true
    if len(login) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, dserr.Err(err)
}

func validateBSPass(pass string) (bool, error) {
    var err error
    var ok bool = true
    if len(pass) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, dserr.Err(err)
}
