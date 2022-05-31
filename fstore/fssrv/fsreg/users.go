/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "ndstore/dscom"
)

const usersSchema = `
    DROP TABLE IF EXISTS users;
    CREATE TABLE IF NOT EXISTS users (
        user_id     INTEGER GENERATED ALWAYS AS IDENTITY (START 1 CYCLE ),
        login       TEXT,
        pass        TEXT,
        state       TEXT,
        role        TEXT
    );
    DROP INDEX IF EXISTS user_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS user_idx
        ON users (login);`


func (reg *Reg) AddUserDescr(login, pass, state, role string) (int64, error) {
    var err error
    request := `
        INSERT INTO users(login, pass, state, role)
        VALUES ($1, $2, $3, $4)
        RETURNING user_id;`
    var userId int64
    err = reg.db.Get(&userId, request, login, pass, state, role)
    if err != nil {
        return userId, err
    }
    return userId, err
}

func (reg *Reg) UserDescrExists(login string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT count(user_id) AS count
        FROM users
        WHERE login = $1
        LIMIT 1;`
    var count int64
    err = reg.db.Get(&count, request, login)
    if err != nil {
        return exists, err
    }
    if count > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) GetUserDescr(login string) (*dscom.UserDescr, error) {
    var err error
    request := `
        SELECT user_id, login, pass, state, role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    user := dscom.NewUserDescr()
    err = reg.db.Get(user, request, login)
    if err != nil {
        return user, err
    }
    return user, err
}

func (reg *Reg) GetUserId(login string) (int64, error) {
    var err error
    request := `
        SELECT user_id
        FROM users
        WHERE login = $1
        LIMIT 1;`
    var userId int64
    err = reg.db.Get(&userId, request, login)
    if err != nil  {
        return userId, err
    }
    return userId, err
}

func (reg *Reg) UpdateUserDescr(login, pass, state, role string) error {
    var err error
    request := `
        UPDATE users
        SET login = $1, pass = $2, state = $3, role = $4
        WHERE login = $1;`
    _, err = reg.db.Exec(request, login, pass, state, role)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) RenewUserDescr(descr *dscom.UserDescr) error {
    var err error
    request := `
        UPDATE users
        SET login = $1, pass = $2, state = $3, role = $4
        WHERE user_id = $5;`
    _, err = reg.db.Exec(request, descr.Login, descr.Pass, descr.State, descr.Role, descr.UserId)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) DeleteUserDescr(login string) error {
    var err error
    request := `
        DELETE FROM users
        WHERE login = $1;`
    _, err = reg.db.Exec(request, login)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListUserDescrs() ([]*dscom.UserDescr, error) {
    var err error
    request := `
        SELECT user_id, login, pass, state, role
        FROM users;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return users, err
    }
    return users, err
}
