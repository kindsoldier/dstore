/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "time"
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
        role        TEXT,
        created_at  INTEGER,
        updated_at  INTEGER
    );
    --- DROP INDEX IF EXISTS fs_user_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS fs_user_idx
        ON fs_users(login);`


func (reg *Reg) AddUserDescr(descr *dscom.UserDescr) (int64, error) {
    var err error
    request := `
        INSERT INTO fs_users(login, pass, state, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING user_id;`
    createdAt := time.Now().Unix()
    updatedAt := createdAt
    var userId int64
    err = reg.db.Get(&userId, request, descr.Login, descr.Pass, descr.State, descr.Role,
                                                                createdAt, updatedAt)
    if err != nil {
        return userId, dserr.Err(err)
    }
    return userId, dserr.Err(err)
}

func (reg *Reg) GetUserDescr(login string) (bool, *dscom.UserDescr, error) {
    var err error
    var exists bool
    request := `
        SELECT user_id, login, pass, state, role, created_at, updated_at
        FROM fs_users
        WHERE login = $1
        LIMIT 1;`
    var user *dscom.UserDescr
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request, login)
    if err != nil {
        return exists, user, dserr.Err(err)
    }
    if len(users) > 0 {
        exists = true
        user = users[0]
        return exists, user, dserr.Err(err)
    }
    return exists, user, dserr.Err(err)
}

func (reg *Reg) UpdateUserDescr(descr *dscom.UserDescr) error {
    var err error
    request := `
        UPDATE fs_users
        SET login = $1, pass = $2, state = $3, role = $4, updated_at = $5
        WHERE user_id = $6;`
    updatedAt := time.Now().Unix()
    _, err = reg.db.Exec(request, descr.Login, descr.Pass, descr.State, descr.Role, updatedAt,
                                                                                descr.UserId)
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
        SELECT user_id, login, pass, state, role, created_at, updated_at
        FROM fs_users;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return users, dserr.Err(err)
    }
    return users, dserr.Err(err)
}
