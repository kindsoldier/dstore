/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsureg

import (
    "errors"
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
)

var ErrorNilRef error = errors.New("db ref is nil")

type Reg struct {
    db *sqlx.DB
}

func NewReg() *Reg {
    var reg Reg
    return &reg
}

func (reg *Reg) OpenDB(dbPath string) error {
    var err error
    db, err := sqlx.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }
    err = db.Ping()
    if err != nil {
        return err
    }
    reg.db = db
    return err
}

func (reg *Reg) CloseDB() error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    reg.db.Close()
    return err
}

func (reg *Reg) MigrateDB() error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    _, err = reg.db.Exec(usersSchema)
    if err != nil {
        return err
    }
    return err
}
