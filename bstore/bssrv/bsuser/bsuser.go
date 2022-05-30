/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsuser

import (
    "errors"
    "ndstore/bstore/bssrv/bsureg"
    "ndstore/bstore/bscom"
)

const UserEnabled string= "enabled"

type Auth struct {
    reg *bsureg.Reg
}

func NewAuth(reg *bsureg.Reg) *Auth {
    var auth Auth
    auth.reg = reg
    return &auth
}


func (auth *Auth) SeedUsers() error {
    var err error
    err = auth.reg.AddUserDescr("admin", "admin", UserEnabled)
    return err
}


func (auth *Auth) AddUser(login, pass string) error {
    var err error
    err = auth.reg.AddUserDescr(login, pass, UserEnabled)
    return err
}

func (auth *Auth) GetUser(login string) (*bscom.UserDescr, bool, error) {
    var err error
    user, exists, err := auth.reg.GetUserDescr(login)
    return user, exists, err
}

func (auth *Auth) CheckUser(login, pass string) (bool, error) {
    var err error
    user, ok, err := auth.reg.GetUserDescr(login)
    if err != nil {
        return ok, err
    }
    if !ok {
        return ok, errors.New("user not exists")
    }
    if pass != user.Pass {
        ok = false
    }
    return ok, err
}


func (auth *Auth) UpdateUser(login, pass string) error {
    var err error
    err = auth.reg.UpdateUserDescr(login, pass, UserEnabled)
    return err
}

func (auth *Auth) ListUsers() ([]*bscom.UserDescr, error) {
    var err error
    users, err := auth.reg.ListUserDescrs()
    return users, err
}

func (auth *Auth) DeleteUser(login string) error {
    var err error
    err = auth.reg.DeleteUserDescr(login)
    return err
}
