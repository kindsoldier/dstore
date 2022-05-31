/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "ndstore/dscom"
)

const StateEnabled  string  = "enabled"

const RoleAdmin string  = "admin"

func (store *Store) AddUser(login, pass string) error {
    var err error
    var ok bool
    ok, err = checkLogin(login)
    if !ok {
        return err
    }
    ok, err = checkPass(pass)
    if !ok {
        return err
    }
    _, err = store.reg.AddUserDescr(login, pass, StateEnabled, RoleAdmin)
    if err != nil {
        return err
    }
    return err
}

func (store *Store) GetUser(login string) (*dscom.UserDescr, error) {
    var err error
    user, err := store.reg.GetUserDescr(login)
    return user,err
}

func (store *Store) CheckUser(login, pass string) (bool, error) {
    var err error
    var ok bool
    user, err := store.reg.GetUserDescr(login)

    if err != nil {
        return ok, err
    }
    if pass == user.Pass {
        ok = true
    }
    return ok, err
}

func (store *Store) UpdateUser(login, pass string) error {
    var err error
    ok, err := checkPass(pass)
    if !ok {
        return err
    }
    err = store.reg.UpdateUserDescr(login, pass, StateEnabled, RoleAdmin)
    return err
}

func (store *Store) ListUsers() ([]*dscom.UserDescr, error) {
    var err error
    users, err := store.reg.ListUserDescrs()
    //for i := range users {
    //    users[i].Pass = "xxxxx"
    //}
    return users, err
}

func (store *Store) DeleteUser(login string) error {
    var err error
    err = store.reg.DeleteUserDescr(login)
    return err
}


func checkLogin(login string) (bool, error) {
    var err error
    var ok bool = true
    if len(login) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, err
}

func checkPass(pass string) (bool, error) {
    var err error
    var ok bool = true
    if len(pass) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, err
}
