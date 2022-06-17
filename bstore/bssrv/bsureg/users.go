/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsureg

import (
    "ndstore/bstore/bscom"
    "ndstore/dserr"
)

const usersSchema = `
    --- DROP TABLE IF EXISTS users;
    CREATE TABLE IF NOT EXISTS users (
        login       TEXT,
        pass        TEXT,
        state       TEXT,
        role        TEXT
    );
    --- DROP INDEX IF EXISTS user_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS user_idx
        ON users (login);`



func (reg *Reg) AddUserDescr(descr *bscom.UserDescr) error {
    var err error
    request := `
        INSERT INTO users(login, pass, state, role)
        VALUES ($1, $2, $3, $4);`
    _, err = reg.db.Exec(request, descr.Login, descr.Pass, descr.State, descr.Role)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) UpdateUserDescr(descr *bscom.UserDescr) error {
    var err error
    request := `
        UPDATE users
        SET pass = $1, state = $2, role = $3
        WHERE login = $4;`
    _, err = reg.db.Exec(request, descr.Pass, descr.State, descr.Role, descr.Login)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) GetUserDescr(login string) (*bscom.UserDescr, error) {
    var err error
    request := `
        SELECT login, pass, state, role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    user := bscom.NewUserDescr()
    err = reg.db.Get(user, request, login)
    if err != nil {
        return user, dserr.Err(err)
    }
    return user, dserr.Err(err)
}

func (reg *Reg) GetUserRole(login string) (string, error) {
    var err error
    request := `
        SELECT role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    var role string
    err = reg.db.Get(&role, request, login)
    if err != nil  {
        return role, dserr.Err(err)
    }
    return role, dserr.Err(err)
}

func (reg *Reg) GetUserState(login string) (string, error) {
    var err error
    request := `
        SELECT state
        FROM users
        WHERE login = $1
        LIMIT 1;`
    var state string
    err = reg.db.Get(&state, request, login)
    if err != nil  {
        return state, dserr.Err(err)
    }
    return state, dserr.Err(err)
}

func (reg *Reg) UserDescrExists(login string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(login) AS count
        FROM users
        WHERE login = $1
        LIMIT 1;`
    var count int64
    err = reg.db.Get(&count, request, login)
    if err != nil {
        return exists, dserr.Err(err)
    }
    if count > 0 {
        exists = true
    }
    return exists, dserr.Err(err)
}

func (reg *Reg) EraseUserDescr(login string) error {
    var err error
    request := `
        DELETE FROM users
        WHERE login = $1;`
    _, err = reg.db.Exec(request, login)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListUserDescrs() ([]*bscom.UserDescr, error) {
    var err error
    request := `
        SELECT login, pass, state, role
        FROM users;`
    users := make([]*bscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return users, dserr.Err(err)
    }
    return users, dserr.Err(err)
}
