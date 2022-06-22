/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "errors"

    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v4/stdlib"
    "ndstore/dserr"

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
    db, err := sqlx.Open("pgx", dbPath)

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
    if reg.db != nil {
        reg.db.Close()
    }
    return dserr.Err(err)
}

func (reg *Reg) MigrateDB() error {
    var err error
    if reg.db == nil {
        return ErrorNilRef
    }
    _, err = reg.db.Exec(blockSchema)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = reg.db.Exec(batchSchema)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = reg.db.Exec(fileSchema)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = reg.db.Exec(entrieSchema)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = reg.db.Exec(bstoreSchema)
    if err != nil {
        return dserr.Err(err)
    }
    _, err = reg.db.Exec(userSchema)
    if err != nil {
        return dserr.Err(err)
    }

    return dserr.Err(err)
}
