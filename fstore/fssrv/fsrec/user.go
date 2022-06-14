/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "ndstore/dscom"
    "ndstore/dserr"
)

const UStateEnabled  string  = "enabled"
const UStateDisabled string  = "disabled"

const URoleAdmin    string  = "admin"
const URoleUser     string  = "user"

const defaultAUser   string  = "admin"
const defaultAPass   string  = "admin"

const defaultUser   string  = "user"
const defaultPass   string  = "user"

func (store *Store) SeedUsers() error {
    var err error
    users, err := store.reg.ListUserDescrs()
    if err != nil {
        return dserr.Err(err)
    }
    if len(users) < 1 {
        _, err = store.reg.AddUserDescr(defaultAUser, defaultAPass, UStateEnabled, URoleAdmin)
        if err != nil {
            return dserr.Err(err)
        }
        _, err = store.reg.AddUserDescr(defaultUser, defaultPass, UStateEnabled, URoleUser)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (store *Store) AddUser(userName, login, pass string) error {
    var err error
    var ok bool

    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        err = errors.New("insufficient rights for adding users")
        return dserr.Err(err)
    }
    ok, err = validateLogin(login)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validatePass(pass)
    if !ok {
        return dserr.Err(err)
    }
    _, err = store.reg.AddUserDescr(login, pass, UStateEnabled, URoleAdmin)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) GetUser(login string) (*dscom.UserDescr, error) {
    var err error
    user, err := store.reg.GetUserDescr(login)
    return user,err
}

func (store *Store) CheckUser(userName, login, pass string) (bool, error) {
    var err error
    var ok bool

    if len(login) == 0 {
        login = userName
    }

    userRole, err := store.reg.GetUserRole(userName)
    if userRole != URoleAdmin && userName != login {
        err = errors.New("insufficient rights for checking other users")
        return ok, dserr.Err(err)
    }

    user, err := store.reg.GetUserDescr(login)

    if err != nil {
        return ok, dserr.Err(err)
    }
    if pass == user.Pass {
        ok = true
    }
    return ok, dserr.Err(err)
}

func (store *Store) UpdateUser(userName, login, newPass, newRole, newState string) error {
    var err error

    // Get current role
    userRole, err := store.reg.GetUserRole(userName)
    if err != nil {
        return dserr.Err(err)
    }

    // Set defaults
    if len(login) < 1 {
        login = userName
    }

    // Rigth control
    if  userName != login && userRole != URoleAdmin {
        err = errors.New("insufficient rights for updating other users")
        return dserr.Err(err)
    }

    // Get old profile and copy to new
    oldUserDescr, err := store.reg.GetUserDescr(login)
    if err != nil {
        return dserr.Err(err)
    }
    newUserDescr := dscom.NewUserDescr()
    newUserDescr.UserId = oldUserDescr.UserId
    newUserDescr.Login  = oldUserDescr.Login
    newUserDescr.Pass   = oldUserDescr.Pass
    newUserDescr.Role   = oldUserDescr.Role
    newUserDescr.State  = oldUserDescr.State

    // Update property if exists
    if len(newPass) > 0 {
        newUserDescr.Pass = newPass
    }
    if len(newRole) > 0 {
        newUserDescr.Role = newRole
    }
    if len(newState) > 0 {
        newUserDescr.State = newState
    }

    // Rigth control
    if newUserDescr.Role != oldUserDescr.Role && userRole != URoleAdmin {
        err = errors.New("insufficient rights for changing role")
        return dserr.Err(err)
    }
    if newUserDescr.State != oldUserDescr.State && userRole != URoleAdmin {
        err = errors.New("insufficient rights for changing state")
        return dserr.Err(err)
    }

    // Validation new property
    var ok bool
    ok, err = validateUState(newUserDescr.State)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateURole(newUserDescr.Role)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validatePass(newUserDescr.Pass)
    if !ok {
        return dserr.Err(err)
    }

    // Update user profile
    err = store.reg.RenewUserDescr(newUserDescr)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) ListUsers(userName string) ([]*dscom.UserDescr, error) {
    var err error
    users := make([]*dscom.UserDescr, 0)
    userRole, err := store.reg.GetUserRole(userName)
    if userRole != URoleAdmin {
        err = errors.New("insufficient rights for listing users")
        return users, dserr.Err(err)
    }
    users, err = store.reg.ListUserDescrs()
    //for i := range users {
    //    users[i].Pass = "xxxxx"
    //}
    if err != nil {
        return users, dserr.Err(err)
    }
    return users, dserr.Err(err)
}

func (store *Store) DeleteUser(userName string, login string) error {
    var err error

    userRole, err := store.reg.GetUserRole(userName)
    if userRole != URoleAdmin {
        err = errors.New("insufficient rights for deleting users")
        return dserr.Err(err)
    }

    err = store.reg.DeleteUserDescr(login)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func validateURole(role string) (bool, error) {
    var err error
    var ok bool = true
    if role == URoleAdmin  {
        return ok, dserr.Err(err)
    }
    if role == URoleUser  {
        return ok, dserr.Err(err)
    }
    err = errors.New("irrelevant role name")
    ok = false
    return ok, dserr.Err(err)
}

func validateUState(state string) (bool, error) {
    var err error
    var ok bool = true
    if state == UStateDisabled  {
        return ok, dserr.Err(err)
    }
    if state == UStateEnabled  {
        return ok, dserr.Err(err)
    }
    err = errors.New("irrelevant state name")
    ok = false
    return ok, dserr.Err(err)
}

func validateLogin(login string) (bool, error) {
    var err error
    var ok bool = true
    if len(login) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, dserr.Err(err)
}

func validatePass(pass string) (bool, error) {
    var err error
    var ok bool = true
    if len(pass) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, dserr.Err(err)
}
