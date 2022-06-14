/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsureg

import (
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
    "ndstore/dserr"
)


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
        return dserr.Err(err)
    }
    err = db.Ping()
    if err != nil {
        return dserr.Err(err)
    }
    reg.db = db
    return dserr.Err(err)
}

func (reg *Reg) CloseDB() error {
    var err error
    reg.db.Close()
    return dserr.Err(err)
}

func (reg *Reg) MigrateDB() error {
    var err error
    _, err = reg.db.Exec(usersSchema)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
