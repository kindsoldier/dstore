/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsbreg

import (
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
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
    reg.db.Close()
    return err
}

func (reg *Reg) MigrateDB() error {
    var err error
    _, err = reg.db.Exec(blockSchema)
    if err != nil {
        return err
    }
    return err
}
