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
        id          INTEGER UNIQUE,
        login       TEXT,
        pass        TEXT,
        state       TEXT
    );
    DROP INDEX IF EXISTS user_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS user_idx
        ON users (id);
    `

func (reg *Reg) AddUserDescr(id int64, login, pass, state string) error {
    var err error
    request := `
        INSERT INTO users(id, login, pass, state)
        VALUES ($1, $2, $3, $4);`
    _, err = reg.db.Exec(request, id, login, pass, state)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) UpdateUserDescr(id int64, login, pass, state string) error {
    var err error
    request := `
        UPDATE users
        SET login = $1, pass = $2, state = $3
        WHERE id = $4;`
    _, err = reg.db.Exec(request, login, pass, state, id)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) RenewUserDescr(descr *dscom.UserDescr) error {
    var err error
    request := `
        UPDATE users
        SET login = $1, pass = $2, state = $3
        WHERE id = $4;`
    _, err = reg.db.Exec(request, descr.Login, descr.Pass, descr.State, descr.Id)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) GetUserDescr(id int64) (*dscom.UserDescr, bool, error) {
    var err error
    var exists bool
    var user *dscom.UserDescr
    request := `
        SELECT id, login, pass, state
        FROM users
        WHERE id = $1
        LIMIT 1;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request, id)
    if err != nil {
        return user, exists, err

    }
    if len(users) > 0 {
        exists = true
        user = users[0]
    }
    return user, exists, err
}

func (reg *Reg) UserDescrExists(id int64) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT id, login, pass, state
        FROM users
        WHERE id = $1
        LIMIT 1;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request, id)
    if err != nil {
        return exists, err
    }
    if len(users) > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) DeleteUserDescr(id int64) error {
    var err error
    request := `
        DELETE FROM users
        WHERE id = $1;`
    _, err = reg.db.Exec(request, id)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListUsersDescr() ([]*dscom.UserDescr, error) {
    var err error
    request := `
        SELECT id, login, pass, state
        FROM users
        WHERE;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return users, err
    }
    return users, err
}
