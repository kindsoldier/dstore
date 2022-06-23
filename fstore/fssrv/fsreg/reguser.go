/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
    "ndstore/dserr"
)

const userSchema = `
    --- DROP TABLE IF EXISTS fs_users;
    CREATE TABLE IF NOT EXISTS fs_users (
        user_id     INTEGER GENERATED ALWAYS AS IDENTITY (START 1 CYCLE ),
        login       TEXT,
        pass        TEXT,
        state       TEXT,
        role        TEXT
    );
    --- DROP INDEX IF EXISTS fs_user_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_user_idx
        ON fs_users(login);`


func (reg *Reg) AddUserDescr(descr *dscom.UserDescr) (int64, error) {
    var err error
    request := `
        INSERT INTO fs_users(login, pass, state, role)
        VALUES ($1, $2, $3, $4)
        RETURNING user_id;`
    var userId int64
    err = reg.db.Get(&userId, request, descr.Login, descr.Pass, descr.State, descr.Role)
    if err != nil {
        return userId, dserr.Err(err)
    }
    return userId, dserr.Err(err)
}

func (reg *Reg) UserDescrExists(login string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(user_id) AS count
        FROM fs_users
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

func (reg *Reg) GetUserDescr(login string) (*dscom.UserDescr, error) {
    var err error
    request := `
        SELECT user_id, login, pass, state, role
        FROM fs_users
        WHERE login = $1
        LIMIT 1;`
    user := dscom.NewUserDescr()
    err = reg.db.Get(user, request, login)
    if err != nil {
        return user, dserr.Err(err)
    }
    return user, dserr.Err(err)
}

func (reg *Reg) GetUserId(login string) (int64, error) {
    var err error
    request := `
        SELECT user_id
        FROM fs_users
        WHERE login = $1
        LIMIT 1;`
    var userId int64
    err = reg.db.Get(&userId, request, login)
    if err != nil  {
        return userId, dserr.Err(err)
    }
    return userId, dserr.Err(err)
}

func (reg *Reg) GetUserRole(login string) (string, error) {
    var err error
    request := `
        SELECT role
        FROM fs_users
        WHERE login = $1
        LIMIT 1;`
    var role string
    err = reg.db.Get(&role, request, login)
    if err != nil  {
        return role, dserr.Err(err)
    }
    return role, dserr.Err(err)
}

func (reg *Reg) UpdateUserDescr(descr *dscom.UserDescr) error {
    var err error
    request := `
        UPDATE fs_users
        SET login = $1, pass = $2, state = $3, role = $4
        WHERE user_id = $5;`
    _, err = reg.db.Exec(request, descr.Login, descr.Pass, descr.State, descr.Role, descr.UserId)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) EraseUserDescr(login string) error {
    var err error
    request := `
        DELETE FROM fs_users
        WHERE login = $1;`
    _, err = reg.db.Exec(request, login)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (reg *Reg) ListUserDescrs() ([]*dscom.UserDescr, error) {
    var err error
    request := `
        SELECT user_id, login, pass, state, role
        FROM fs_users;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return users, dserr.Err(err)
    }
    return users, dserr.Err(err)
}
