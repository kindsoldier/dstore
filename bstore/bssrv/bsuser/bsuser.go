/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsuser

import (
    "errors"
    "fmt"
    "ndstore/bstore/bssrv/bsureg"
    "ndstore/bstore/bscom"
    "ndstore/dserr"
)

const UStateEnabled  string  = "enabled"
const UStateDisabled string  = "disabled"

const URoleAdmin    string  = "admin"
const URoleUser     string  = "user"

const defaultAUser  string  = "admin"
const defaultAPass  string  = "admin"

const defaultUser   string  = "user"
const defaultPass   string  = "user"

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
    users, err := auth.reg.ListUserDescrs()
    if err != nil {
        return dserr.Err(err)
    }
    if len(users) < 1 {
        err = auth.reg.AddUserDescr(defaultAUser, defaultAPass, UStateEnabled, URoleAdmin)
        if err != nil {
            return dserr.Err(err)
        }
        err = auth.reg.AddUserDescr(defaultUser, defaultPass, UStateEnabled, URoleUser)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (auth *Auth) AddUser(userName, login, pass, role, state string) error {
    var err error
    var ok bool

    // Check user role
    userRole, err := auth.reg.GetUserRole(userName)
    if userRole != URoleAdmin {
        err = errors.New("insufficient rights for adding users")
        return dserr.Err(err)
    }

    // Set defaults
    if len(role) < 1 {
        role = URoleUser
    }
    if len(state) < 1 {
        state = UStateEnabled
    }

    // Validate parameters
    ok, err = validateUState(state)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateURole(role)
    if !ok {
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
    err = auth.reg.AddUserDescr(login, pass, state, role)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (auth *Auth) UserExists(login string) (bool, error) {
    var err error
    exists, err := auth.reg.UserDescrExists(login)
    return exists, dserr.Err(err)
}


func (auth *Auth) GetUser(login string) (*bscom.UserDescr, error) {
    var err error
    user, err := auth.reg.GetUserDescr(login)
    return user, dserr.Err(err)
}

func (auth *Auth) CheckUser(userName, login, pass string) (bool, error) {
    var err error
    var ok bool

    if len(login) == 0 {
        login = userName
    }

    userRole, err := auth.reg.GetUserRole(userName)
    if userRole != URoleAdmin && userName != login {
        err = errors.New("insufficient rights for checking other users")
        return ok, dserr.Err(err)
    }

    user, err := auth.reg.GetUserDescr(login)

    if err != nil {
        return ok, dserr.Err(err)
    }
    if pass == user.Pass {
        ok = true
    }
    return ok, dserr.Err(err)
}

func (auth *Auth) UpdateUser(userName, login, newPass, newRole, newState string) error {
    var err error

    // Get current role
    userRole, err := auth.reg.GetUserRole(userName)
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
    oldUserDescr, err := auth.reg.GetUserDescr(login)
    if err != nil {
        return dserr.Err(err)
    }
    newUserDescr := bscom.NewUserDescr()
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
    err = auth.reg.RenewUserDescr(newUserDescr)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (auth *Auth) ListUsers(userName string) ([]*bscom.UserDescr, error) {
    var err error
    users := make([]*bscom.UserDescr, 0)
    userRole, err := auth.reg.GetUserRole(userName)
    if userRole != URoleAdmin {
        err = errors.New("insufficient rights for listing users")
        return users, dserr.Err(err)
    }
    users, err = auth.reg.ListUserDescrs()
    //for i := range users {
    //    users[i].Pass = "xxxxx"
    //}
    if err != nil {
        return users, dserr.Err(err)
    }
    return users, dserr.Err(err)
}

func (auth *Auth) DeleteUser(userName string, login string) error {
    var err error

    userRole, err := auth.reg.GetUserRole(userName)
    if userRole != URoleAdmin {
        err = errors.New("insufficient rights for deleting users")
        return dserr.Err(err)
    }

    err = auth.reg.DeleteUserDescr(login)
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
    err = fmt.Errorf("irrelevant role name: %s", role)
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
    err = fmt.Errorf("irrelevant state name: %s", state)
    ok = false
    return ok, dserr.Err(err)
}

func validateLogin(login string) (bool, error) {
    var err error
    var ok bool = true
    if len(login) < 1 {
        ok = false
        err = errors.New("zero len of login")
    }
    return ok, dserr.Err(err)
}

func validatePass(pass string) (bool, error) {
    var err error
    var ok bool = true
    if len(pass) < 1 {
        ok = false
        err = errors.New("zero len of password")
    }
    return ok, dserr.Err(err)
}
