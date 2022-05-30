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
        state       TEXT,
        role        TEXT
    );
    DROP INDEX IF EXISTS user1_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS user1_idx
        ON users (id);
    DROP INDEX IF EXISTS user2_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS user2_idx
        ON users (login);
    DROP INDEX IF EXISTS user3_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS user3_idx
        ON users (id, login);`


func (reg *Reg) GetNewUserId() (int64, error) {
    var err error
    var userId int64
    request := `
        SELECT id
        FROM users
        ORDER BY id DESC
        LIMIT 1;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return userId, err
    }
    if len(users) > 0 {
        userId = users[0].Id + 1
    }
    return userId, err
}

func (reg *Reg) AddUserDescr(id int64, login, pass, state, role string) error {
    var err error
    request := `
        INSERT INTO users(id, login, pass, state, role)
        VALUES ($1, $2, $3, $4, $5);`
    _, err = reg.db.Exec(request, id, login, pass, state, role)
    if err != nil {
        return err
    }
    return err
}


func (reg *Reg) GetUserDescr(login string) (*dscom.UserDescr, bool, error) {
    var err error
    var exists bool
    var user *dscom.UserDescr
    request := `
        SELECT id, login, pass, state, role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request, login)
    if err != nil {
        return user, exists, err

    }
    if len(users) > 0 {
        exists = true
        user = users[0]
    }
    return user, exists, err
}

func (reg *Reg) GetUserId(login string) (int64, bool, error) {
    var err error
    var exists bool
    var userId int64
    request := `
        SELECT id, login, pass, state, role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request, login)
    if err != nil {
        return userId, exists, err

    }
    if len(users) > 0 {
        exists = true
        userId = users[0].Id
    }
    return userId, exists, err
}

func (reg *Reg) UpdateUserDescr(login, pass, state, role string) error {
    var err error
    request := `
        UPDATE users
        SET pass = $1, state = $2, role = $3
        WHERE login = $4;`
    _, err = reg.db.Exec(request, pass, state, role, login)
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
        WHERE id = $5;`
    _, err = reg.db.Exec(request, descr.Login, descr.Pass, descr.State, descr.Role, descr.Id)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) UserDescrExists(login string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT id, login, pass, state, role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request, login)
    if err != nil {
        return exists, err
    }
    if len(users) > 0 {
        exists = true
    }
    return exists, err
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
        SELECT id, login, pass, state, role
        FROM users;`
    users := make([]*dscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return users, err
    }
    return users, err
}
