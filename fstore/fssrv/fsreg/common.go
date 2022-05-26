/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package fsreg

import (
    "errors"

    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v4/stdlib"
)

const entriesSchema = `
    DROP TABLE IF EXISTS entries;
    CREATE TABLE IF NOT EXISTS entries (
        dir_path   TEXT,
        file_name  TEXT,
        file_id    INTEGER
    );
    `

const filesSchema = `
    DROP TABLE IF EXISTS files;
    CREATE TABLE IF NOT EXISTS files (
        file_id     INTEGER,
        batch_count INTEGER,
        batch_size  INTEGER,
        block_size  INTEGER,
        file_size   INTEGER
    );
    DROP INDEX IF EXISTS file_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS file_idx
        ON files (file_id);

    DROP TABLE IF EXISTS batchs;
    CREATE TABLE IF NOT EXISTS batchs (
        file_id     INTEGER,
        batch_id    INTEGER,
        batch_size  INTEGER,
        block_size  INTEGER
    );
    DROP INDEX IF EXISTS batch_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS batch_idx
        ON batchs (file_id, batch_id);


    DROP TABLE IF EXISTS blocks;
    CREATE TABLE IF NOT EXISTS blocks (
        file_id     INTEGER,
        batch_id    INTEGER,
        block_id    INTEGER,
        block_size  INTEGER,
        file_path   TEXT,
        hash_init   TEXT,
        hash_sum    TEXT,
        data_size  INTEGER
    );
    DROP INDEX IF EXISTS block_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS block_idx
        ON blocks (file_id, batch_id, block_id);
    `

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
    _, err = reg.db.Exec(filesSchema)
    if err != nil {
        return err
    }

    _, err = reg.db.Exec(entriesSchema)
    if err != nil {
        return err
    }

    return err
}
